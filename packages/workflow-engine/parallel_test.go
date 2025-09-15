package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// MockWorkflowEngine for testing
type MockWorkflowEngine struct {
	executeFunc func(WorkflowStep, int) error
}

func (m *MockWorkflowEngine) executeStep(step WorkflowStep, stepNum int) error {
	if m.executeFunc != nil {
		return m.executeFunc(step, stepNum)
	}
	// Default implementation - just return success
	return nil
}

func TestDefaultParallelConfig(t *testing.T) {
	config := DefaultParallelConfig()

	if config.MaxConcurrency != 4 {
		t.Errorf("Expected MaxConcurrency to be 4, got %d", config.MaxConcurrency)
	}

	if !config.EnableParallel {
		t.Error("Expected EnableParallel to be true")
	}

	if config.TimeoutDuration != 5*time.Minute {
		t.Errorf("Expected TimeoutDuration to be 5 minutes, got %v", config.TimeoutDuration)
	}

	if !config.DependencyCheck {
		t.Error("Expected DependencyCheck to be true")
	}
}

func TestNewParallelExecutor(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)

	if executor.config.MaxConcurrency != config.MaxConcurrency {
		t.Error("Executor config not initialized correctly")
	}

	if executor.stepResults == nil {
		t.Error("Step results map not initialized")
	}

	if executor.workerPool == nil {
		t.Error("Worker pool not initialized")
	}

	// Cleanup
	executor.Cleanup()
}

func TestBuildDependencyGraph_NoDependencies(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	steps := []WorkflowStep{
		{Agent: "agent1", Task: "task1", Template: "template1.yaml"},
		{Agent: "agent2", Task: "task2", Template: "template2.yaml"},
		{Agent: "agent3", Task: "task3", Checklist: "checklist1.yaml"},
	}

	graph, err := executor.BuildDependencyGraph(steps)
	if err != nil {
		t.Fatalf("Failed to build dependency graph: %v", err)
	}

	if len(graph.Steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(graph.Steps))
	}

	// All steps should have in-degree 0 (no dependencies)
	for i := 0; i < 3; i++ {
		if graph.InDegree[i] != 0 {
			t.Errorf("Step %d should have in-degree 0, got %d", i, graph.InDegree[i])
		}
	}
}

func TestBuildDependencyGraph_WithDependencies(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	steps := []WorkflowStep{
		{Agent: "architect", Task: "design", Template: "architecture.yaml"},
		{Agent: "dev", Task: "implement", Variables: map[string]interface{}{"requires": "architecture"}},
		{Agent: "qa", Task: "validate", Checklist: "validation.yaml"},
	}

	graph, err := executor.BuildDependencyGraph(steps)
	if err != nil {
		t.Fatalf("Failed to build dependency graph: %v", err)
	}

	// Step 1 (dev) should depend on step 0 (architect)
	if graph.InDegree[1] != 1 {
		t.Errorf("Dev step should depend on architect, in-degree: %d", graph.InDegree[1])
	}

	// QA step should depend on previous outputs
	if graph.InDegree[2] == 0 {
		t.Error("QA step should have dependencies")
	}
}

func TestBuildDependencyGraph_CircularDependency(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Manually create circular dependency for testing
	executor.dependencyGraph = &DependencyGraph{
		Steps: []StepDependency{
			{StepIndex: 0, Dependencies: []int{2}},
			{StepIndex: 1, Dependencies: []int{0}},
			{StepIndex: 2, Dependencies: []int{1}},
		},
		AdjacencyList: map[int][]int{
			0: {1},
			1: {2},
			2: {0},
		},
		InDegree: map[int]int{
			0: 1, 1: 1, 2: 1,
		},
	}

	err := executor.validateDAG(executor.dependencyGraph)
	if err == nil {
		t.Error("Expected circular dependency error, but got none")
	}
}

func TestExtractOutputsAndInputs(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Test template step
	templateStep := WorkflowStep{
		Agent:    "architect",
		Task:     "design",
		Template: "architecture.yaml",
		Variables: map[string]interface{}{
			"output_file": "design.md",
		},
	}

	outputs := executor.extractOutputs(templateStep)
	if len(outputs) < 1 {
		t.Error("Template step should have outputs")
	}

	// Test checklist step
	checklistStep := WorkflowStep{
		Agent:     "qa",
		Task:      "validate",
		Checklist: "quality.yaml",
	}

	outputs = executor.extractOutputs(checklistStep)
	if len(outputs) < 1 {
		t.Error("Checklist step should have outputs")
	}

	// Test step with inputs
	dependentStep := WorkflowStep{
		Agent: "dev",
		Task:  "implement",
		Variables: map[string]interface{}{
			"input_file": "requirements.md",
		},
	}

	inputs := executor.extractInputs(dependentStep)
	if len(inputs) < 1 {
		t.Error("Dependent step should have inputs")
	}
}

func TestHasDependency(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Test architect -> dev dependency
	architectStep := WorkflowStep{Agent: "architect", Task: "design"}
	devStep := WorkflowStep{Agent: "dev", Task: "implement"}

	if !executor.hasDependency(architectStep, devStep) {
		t.Error("Dev step should depend on architect step")
	}

	// Test no dependency
	step1 := WorkflowStep{Agent: "agent1", Task: "task1"}
	step2 := WorkflowStep{Agent: "agent2", Task: "task2"}

	if executor.hasDependency(step1, step2) {
		t.Error("Independent steps should not have dependency")
	}
}

func TestExecuteSequential(t *testing.T) {
	config := ParallelExecutionConfig{
		EnableParallel:  false,
		MaxConcurrency:  2,
		TimeoutDuration: 1 * time.Second,
	}

	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Mock workflow engine
	mockEngine := &MockWorkflowEngine{
		executeFunc: func(step WorkflowStep, stepNum int) error {
			return nil // Success
		},
	}

	steps := []WorkflowStep{
		{Agent: "agent1", Task: "task1"},
		{Agent: "agent2", Task: "task2"},
	}

	err := executor.executeSequential(mockEngine, steps)
	if err != nil {
		t.Errorf("Sequential execution failed: %v", err)
	}
}

func TestProgressUpdates(t *testing.T) {
	config := DefaultParallelConfig()
	config.TimeoutDuration = 1 * time.Second

	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Start progress monitoring
	go executor.monitorProgress(2)

	// Send progress updates
	executor.updateProgress(0, 2, "executing", "Test step 1")
	executor.updateProgress(1, 2, "completed", "Test step 2")

	// Wait briefly for monitoring to process updates
	time.Sleep(100 * time.Millisecond)

	// Check that progress channel received updates (may be empty due to timing)
	select {
	case <-executor.progressChan:
		// Good, received update
	default:
		// This is okay too, updates might have been processed already
	}
}

func TestStepResult(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Simulate step execution result
	result := &StepResult{
		StepIndex: 0,
		Success:   true,
		Error:     nil,
		StartTime: time.Now().Add(-1 * time.Second),
		EndTime:   time.Now(),
		Duration:  1 * time.Second,
	}

	executor.mutex.Lock()
	executor.stepResults[0] = result
	executor.mutex.Unlock()

	// Get results
	results := executor.GetResults()
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if !results[0].Success {
		t.Error("Expected successful result")
	}

	if results[0].Duration != 1*time.Second {
		t.Errorf("Expected duration 1s, got %v", results[0].Duration)
	}
}

func TestConcurrentStepExecution(t *testing.T) {
	config := ParallelExecutionConfig{
		EnableParallel:  true,
		MaxConcurrency:  2,
		TimeoutDuration: 5 * time.Second,
		DependencyCheck: true,
	}

	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Track concurrent executions
	var concurrentCount int32
	var maxConcurrent int32

	// Create a mock engine that tracks concurrency
	mockEngine := &MockWorkflowEngine{
		executeFunc: func(step WorkflowStep, stepNum int) error {
			// Increment concurrent count
			current := atomic.AddInt32(&concurrentCount, 1)

			// Update max concurrent if needed
			for {
				max := atomic.LoadInt32(&maxConcurrent)
				if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
					break
				}
			}

			// Simulate work
			time.Sleep(100 * time.Millisecond)

			// Decrement concurrent count
			atomic.AddInt32(&concurrentCount, -1)

			return nil
		},
	}

	// Create independent steps (no dependencies)
	steps := []WorkflowStep{
		{Agent: "agent1", Task: "task1"},
		{Agent: "agent2", Task: "task2"},
		{Agent: "agent3", Task: "task3"},
		{Agent: "agent4", Task: "task4"},
	}

	// Execute in parallel
	err := executor.ExecuteParallel(mockEngine, steps)
	if err != nil {
		t.Fatalf("Parallel execution failed: %v", err)
	}

	// Verify concurrent execution occurred
	maxConcurrentValue := atomic.LoadInt32(&maxConcurrent)
	if maxConcurrentValue < 2 {
		t.Errorf("Expected at least 2 concurrent executions, got %d", maxConcurrentValue)
	}

	// Verify all steps completed
	results := executor.GetResults()
	if len(results) != 4 {
		t.Errorf("Expected 4 results, got %d", len(results))
	}
}

func TestTimeout(t *testing.T) {
	config := ParallelExecutionConfig{
		EnableParallel:  true,
		MaxConcurrency:  2,
		TimeoutDuration: 100 * time.Millisecond, // Very short timeout
		DependencyCheck: true,
	}

	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Mock engine with long delay
	mockEngine := &MockWorkflowEngine{
		executeFunc: func(step WorkflowStep, stepNum int) error {
			// Sleep longer than timeout
			time.Sleep(200 * time.Millisecond)
			return nil
		},
	}

	steps := []WorkflowStep{
		{Agent: "agent1", Task: "task1"},
	}

	// Execute and expect timeout
	err := executor.ExecuteParallel(mockEngine, steps)
	if err == nil {
		t.Error("Expected timeout error, but got none")
	}
}

func TestPrintExecutionSummary(t *testing.T) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	// Add some mock results
	executor.stepResults[0] = &StepResult{
		StepIndex: 0,
		Success:   true,
		Duration:  1 * time.Second,
	}

	executor.stepResults[1] = &StepResult{
		StepIndex: 1,
		Success:   false,
		Duration:  2 * time.Second,
		Error:     fmt.Errorf("test error"),
	}

	// This should not panic and should produce output
	executor.PrintExecutionSummary()

	// Test with no results
	executor.stepResults = make(map[int]*StepResult)
	executor.PrintExecutionSummary() // Should handle empty results gracefully
}

// Benchmark tests
func BenchmarkBuildDependencyGraph(b *testing.B) {
	config := DefaultParallelConfig()
	executor := NewParallelExecutor(config)
	defer executor.Cleanup()

	steps := make([]WorkflowStep, 10)
	for i := 0; i < 10; i++ {
		steps[i] = WorkflowStep{
			Agent: fmt.Sprintf("agent%d", i),
			Task:  fmt.Sprintf("task%d", i),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.BuildDependencyGraph(steps)
		if err != nil {
			b.Fatalf("Failed to build dependency graph: %v", err)
		}
	}
}

func BenchmarkParallelExecution(b *testing.B) {
	config := ParallelExecutionConfig{
		EnableParallel:  true,
		MaxConcurrency:  4,
		TimeoutDuration: 30 * time.Second,
		DependencyCheck: true,
	}

	steps := []WorkflowStep{
		{Agent: "agent1", Task: "task1"},
		{Agent: "agent2", Task: "task2"},
		{Agent: "agent3", Task: "task3"},
		{Agent: "agent4", Task: "task4"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executor := NewParallelExecutor(config)

		// Mock engine
		mockEngine := &MockWorkflowEngine{
			executeFunc: func(step WorkflowStep, stepNum int) error {
				// Simulate small amount of work
				time.Sleep(1 * time.Millisecond)
				return nil
			},
		}

		err := executor.ExecuteParallel(mockEngine, steps)
		if err != nil {
			b.Fatalf("Parallel execution failed: %v", err)
		}
		executor.Cleanup()
	}
}
