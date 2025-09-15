package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ParallelExecutionConfig holds configuration for parallel execution
type ParallelExecutionConfig struct {
	MaxConcurrency  int           `yaml:"max_concurrency,omitempty"`
	EnableParallel  bool          `yaml:"enable_parallel,omitempty"`
	TimeoutDuration time.Duration `yaml:"timeout_duration,omitempty"`
	DependencyCheck bool          `yaml:"dependency_check,omitempty"`
}

// DefaultParallelConfig returns sensible defaults
func DefaultParallelConfig() ParallelExecutionConfig {
	return ParallelExecutionConfig{
		MaxConcurrency:  4,
		EnableParallel:  true,
		TimeoutDuration: 5 * time.Minute,
		DependencyCheck: true,
	}
}

// StepDependency represents dependencies between workflow steps
type StepDependency struct {
	StepIndex    int      `yaml:"step_index"`
	Dependencies []int    `yaml:"dependencies,omitempty"`
	Outputs      []string `yaml:"outputs,omitempty"`
	Inputs       []string `yaml:"inputs,omitempty"`
}

// DependencyGraph represents the workflow step dependency graph
type DependencyGraph struct {
	Steps         []StepDependency
	AdjacencyList map[int][]int
	InDegree      map[int]int
}

// StepResult holds the result of executing a workflow step
type StepResult struct {
	StepIndex int
	Success   bool
	Error     error
	Output    interface{}
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// ParallelExecutor manages parallel execution of workflow steps
type ParallelExecutor struct {
	config          ParallelExecutionConfig
	dependencyGraph *DependencyGraph
	stepResults     map[int]*StepResult
	workerPool      chan struct{}
	resultChan      chan *StepResult
	errorChan       chan error
	progressChan    chan ProgressUpdate
	mutex           sync.RWMutex
	wg              sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
}

// ProgressUpdate represents real-time progress information
type ProgressUpdate struct {
	StepIndex      int
	TotalSteps     int
	CompletedSteps int
	Status         string
	Message        string
	Timestamp      time.Time
}

// NewParallelExecutor creates a new parallel execution engine
func NewParallelExecutor(config ParallelExecutionConfig) *ParallelExecutor {
	ctx, cancel := context.WithTimeout(context.Background(), config.TimeoutDuration)

	return &ParallelExecutor{
		config:       config,
		stepResults:  make(map[int]*StepResult),
		workerPool:   make(chan struct{}, config.MaxConcurrency),
		resultChan:   make(chan *StepResult, 100),
		errorChan:    make(chan error, 100),
		progressChan: make(chan ProgressUpdate, 100),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// BuildDependencyGraph analyzes workflow steps and builds dependency graph
func (pe *ParallelExecutor) BuildDependencyGraph(steps []WorkflowStep) (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Steps:         make([]StepDependency, len(steps)),
		AdjacencyList: make(map[int][]int),
		InDegree:      make(map[int]int),
	}

	// Initialize step dependencies
	for i, step := range steps {
		stepDep := StepDependency{
			StepIndex:    i,
			Dependencies: []int{},
			Outputs:      pe.extractOutputs(step),
			Inputs:       pe.extractInputs(step),
		}

		// Analyze dependencies based on inputs/outputs
		for j := 0; j < i; j++ {
			if pe.hasDependency(steps[j], step) {
				stepDep.Dependencies = append(stepDep.Dependencies, j)
				graph.AdjacencyList[j] = append(graph.AdjacencyList[j], i)
				graph.InDegree[i]++
			}
		}

		graph.Steps[i] = stepDep

		// Initialize in-degree for steps with no dependencies
		if _, exists := graph.InDegree[i]; !exists {
			graph.InDegree[i] = 0
		}
	}

	// Validate graph (check for cycles)
	if err := pe.validateDAG(graph); err != nil {
		return nil, fmt.Errorf("dependency graph validation failed: %v", err)
	}

	pe.dependencyGraph = graph
	return graph, nil
}

// extractOutputs identifies potential outputs from a workflow step
func (pe *ParallelExecutor) extractOutputs(step WorkflowStep) []string {
	outputs := []string{}

	// Template steps generate files
	if step.Template != "" {
		outputs = append(outputs, "template_output")
	}

	// Checklist steps generate reports
	if step.Checklist != "" {
		outputs = append(outputs, "checklist_report")
	}

	// Extract from variables
	for key := range step.Variables {
		if key == "output_file" || key == "generates" {
			outputs = append(outputs, fmt.Sprintf("var_%s", key))
		}
	}

	return outputs
}

// extractInputs identifies potential inputs for a workflow step
func (pe *ParallelExecutor) extractInputs(step WorkflowStep) []string {
	inputs := []string{}

	// Steps that require previous outputs
	if step.Agent == "qa" || step.Task == "validate" {
		inputs = append(inputs, "template_output", "checklist_report")
	}

	// Extract from variables
	for key := range step.Variables {
		if key == "input_file" || key == "requires" {
			inputs = append(inputs, fmt.Sprintf("var_%s", key))
		}
	}

	return inputs
}

// hasDependency checks if step2 depends on step1
func (pe *ParallelExecutor) hasDependency(step1, step2 WorkflowStep) bool {
	outputs1 := pe.extractOutputs(step1)
	inputs2 := pe.extractInputs(step2)

	// Check for overlap between outputs and inputs
	for _, output := range outputs1 {
		for _, input := range inputs2 {
			if output == input {
				return true
			}
		}
	}

	// Conservative dependency for certain agent combinations
	if step1.Agent == "architect" && step2.Agent == "dev" {
		return true
	}

	if step1.Agent == "po" && (step2.Agent == "sm" || step2.Agent == "dev") {
		return true
	}

	return false
}

// validateDAG checks for cycles in the dependency graph
func (pe *ParallelExecutor) validateDAG(graph *DependencyGraph) error {
	visited := make(map[int]bool)
	recStack := make(map[int]bool)

	for i := 0; i < len(graph.Steps); i++ {
		if !visited[i] {
			if pe.hasCycle(graph, i, visited, recStack) {
				return fmt.Errorf("circular dependency detected involving step %d", i)
			}
		}
	}

	return nil
}

// hasCycle performs DFS to detect cycles
func (pe *ParallelExecutor) hasCycle(graph *DependencyGraph, node int, visited, recStack map[int]bool) bool {
	visited[node] = true
	recStack[node] = true

	for _, neighbor := range graph.AdjacencyList[node] {
		if !visited[neighbor] {
			if pe.hasCycle(graph, neighbor, visited, recStack) {
				return true
			}
		} else if recStack[neighbor] {
			return true
		}
	}

	recStack[node] = false
	return false
}

// ExecuteParallel executes workflow steps in parallel based on dependency graph
func (pe *ParallelExecutor) ExecuteParallel(engine StepExecutor, steps []WorkflowStep) error {
	if !pe.config.EnableParallel {
		return pe.executeSequential(engine, steps)
	}

	// Build dependency graph
	graph, err := pe.BuildDependencyGraph(steps)
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %v", err)
	}

	fmt.Printf("📊 Dependency Analysis Complete:\n")
	fmt.Printf("   🔗 Total Steps: %d\n", len(steps))
	fmt.Printf("   ⚡ Parallel Steps: %d\n", pe.countParallelSteps(graph))
	fmt.Printf("   🎯 Max Concurrency: %d\n", pe.config.MaxConcurrency)

	// Start progress monitoring
	go pe.monitorProgress(len(steps))

	// Execute steps using topological sort
	return pe.executeTopological(engine, steps, graph)
}

// countParallelSteps counts how many steps can run in parallel
func (pe *ParallelExecutor) countParallelSteps(graph *DependencyGraph) int {
	parallelCount := 0
	for _, inDegree := range graph.InDegree {
		if inDegree == 0 {
			parallelCount++
		}
	}
	return parallelCount
}

// executeSequential falls back to sequential execution
func (pe *ParallelExecutor) executeSequential(engine StepExecutor, steps []WorkflowStep) error {
	fmt.Printf("🔄 Sequential Execution Mode (parallel disabled)\n")

	for i, step := range steps {
		select {
		case <-pe.ctx.Done():
			return fmt.Errorf("execution timeout or cancelled")
		default:
			pe.updateProgress(i, len(steps), "executing", fmt.Sprintf("Step %d: %s", i+1, step.Task))

			if err := engine.executeStep(step, i+1); err != nil {
				return fmt.Errorf("step %d failed: %v", i+1, err)
			}
		}
	}

	return nil
}

// executeTopological executes steps using topological sorting for parallel execution
func (pe *ParallelExecutor) executeTopological(engine StepExecutor, steps []WorkflowStep, graph *DependencyGraph) error {
	ready := make([]int, 0)
	inDegree := make(map[int]int)

	// Initialize in-degree tracking
	for k, v := range graph.InDegree {
		inDegree[k] = v
		if v == 0 {
			ready = append(ready, k)
		}
	}

	completed := 0

	// Process steps level by level
	for len(ready) > 0 || completed < len(steps) {
		select {
		case <-pe.ctx.Done():
			pe.cancel()
			return fmt.Errorf("execution timeout or cancelled")
		default:
			// Execute all ready steps in parallel
			if len(ready) > 0 {
				currentBatch := make([]int, len(ready))
				copy(currentBatch, ready)
				ready = ready[:0]

				// Execute batch in parallel
				if err := pe.executeBatch(engine, steps, currentBatch); err != nil {
					return err
				}

				// Update dependencies for completed steps
				for _, stepIndex := range currentBatch {
					completed++

					// Reduce in-degree for dependent steps
					for _, dependent := range graph.AdjacencyList[stepIndex] {
						inDegree[dependent]--
						if inDegree[dependent] == 0 {
							ready = append(ready, dependent)
						}
					}
				}
			} else if completed < len(steps) {
				// Wait for some steps to complete
				time.Sleep(100 * time.Millisecond)

				// Check for deadlock
				allBlocked := true
				for i := 0; i < len(steps); i++ {
					if result, exists := pe.stepResults[i]; !exists || result == nil {
						if inDegree[i] == 0 {
							ready = append(ready, i)
							allBlocked = false
						}
					}
				}

				if allBlocked && len(ready) == 0 {
					return fmt.Errorf("execution deadlock detected - no steps can proceed")
				}
			}
		}
	}

	return nil
}

// executeBatch executes a batch of ready steps in parallel
func (pe *ParallelExecutor) executeBatch(engine StepExecutor, steps []WorkflowStep, batch []int) error {
	batchSize := len(batch)
	if batchSize == 0 {
		return nil
	}

	fmt.Printf("⚡ Executing batch of %d parallel steps: %v\n", batchSize, batch)

	// Use worker pool to limit concurrency
	for _, stepIndex := range batch {
		pe.wg.Add(1)
		go pe.executeStepWorker(engine, steps[stepIndex], stepIndex)
	}

	// Wait for all steps in batch to complete with timeout check
	done := make(chan struct{})
	go func() {
		pe.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All steps completed
	case <-pe.ctx.Done():
		// Timeout occurred
		return fmt.Errorf("batch execution timeout")
	}

	// Check for errors in batch
	for _, stepIndex := range batch {
		if result, exists := pe.stepResults[stepIndex]; exists && result.Error != nil {
			return fmt.Errorf("step %d failed: %v", stepIndex+1, result.Error)
		}
	}

	return nil
}

// StepExecutor interface for testing
type StepExecutor interface {
	executeStep(step WorkflowStep, stepIndex int) error
}

// executeStepWorker executes a single step in a goroutine
func (pe *ParallelExecutor) executeStepWorker(engine StepExecutor, step WorkflowStep, stepIndex int) {
	defer pe.wg.Done()

	// Check for cancellation before starting
	select {
	case <-pe.ctx.Done():
		pe.mutex.Lock()
		pe.stepResults[stepIndex] = &StepResult{
			StepIndex: stepIndex,
			Success:   false,
			Error:     fmt.Errorf("execution cancelled"),
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Duration:  0,
		}
		pe.mutex.Unlock()
		return
	default:
	}

	// Acquire worker slot
	pe.workerPool <- struct{}{}
	defer func() { <-pe.workerPool }()

	startTime := time.Now()

	pe.updateProgress(stepIndex, -1, "executing", fmt.Sprintf("Step %d: %s", stepIndex+1, step.Task))

	// Execute step with error handling
	err := engine.executeStep(step, stepIndex+1)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// Store result
	result := &StepResult{
		StepIndex: stepIndex,
		Success:   err == nil,
		Error:     err,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration,
	}

	pe.mutex.Lock()
	pe.stepResults[stepIndex] = result
	pe.mutex.Unlock()

	// Report result
	if err != nil {
		pe.updateProgress(stepIndex, -1, "failed", fmt.Sprintf("Step %d failed: %v", stepIndex+1, err))
	} else {
		pe.updateProgress(stepIndex, -1, "completed", fmt.Sprintf("Step %d completed in %v", stepIndex+1, duration))
	}
}

// updateProgress sends progress updates
func (pe *ParallelExecutor) updateProgress(stepIndex, totalSteps int, status, message string) {
	pe.mutex.RLock()
	completedCount := len(pe.stepResults)
	pe.mutex.RUnlock()

	update := ProgressUpdate{
		StepIndex:      stepIndex,
		TotalSteps:     totalSteps,
		CompletedSteps: completedCount,
		Status:         status,
		Message:        message,
		Timestamp:      time.Now(),
	}

	select {
	case pe.progressChan <- update:
	default:
		// Channel full, skip update
	}
}

// monitorProgress displays real-time progress updates
func (pe *ParallelExecutor) monitorProgress(totalSteps int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-pe.ctx.Done():
			return
		case update := <-pe.progressChan:
			fmt.Printf("   📊 Progress: [%d/%d] %s - %s\n",
				update.CompletedSteps, totalSteps, update.Status, update.Message)
		case <-ticker.C:
			pe.mutex.RLock()
			completed := len(pe.stepResults)
			pe.mutex.RUnlock()

			if completed >= totalSteps {
				return
			}
		}
	}
}

// GetResults returns the execution results for all steps
func (pe *ParallelExecutor) GetResults() map[int]*StepResult {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()

	results := make(map[int]*StepResult)
	for k, v := range pe.stepResults {
		results[k] = v
	}

	return results
}

// Cleanup releases resources
func (pe *ParallelExecutor) Cleanup() {
	pe.cancel()
	close(pe.resultChan)
	close(pe.errorChan)
	close(pe.progressChan)
}

// PrintExecutionSummary displays execution statistics
func (pe *ParallelExecutor) PrintExecutionSummary() {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()

	totalSteps := len(pe.stepResults)
	if totalSteps == 0 {
		return
	}

	successCount := 0
	totalDuration := time.Duration(0)
	var minDuration, maxDuration time.Duration

	for _, result := range pe.stepResults {
		if result.Success {
			successCount++
		}

		totalDuration += result.Duration

		if minDuration == 0 || result.Duration < minDuration {
			minDuration = result.Duration
		}

		if result.Duration > maxDuration {
			maxDuration = result.Duration
		}
	}

	avgDuration := totalDuration / time.Duration(totalSteps)

	fmt.Printf("\n📈 Parallel Execution Summary:\n")
	fmt.Printf("   ✅ Success Rate: %d/%d (%.1f%%)\n",
		successCount, totalSteps, float64(successCount)/float64(totalSteps)*100)
	fmt.Printf("   ⏱️  Total Duration: %v\n", totalDuration)
	fmt.Printf("   📊 Average Duration: %v\n", avgDuration)
	fmt.Printf("   ⚡ Fastest Step: %v\n", minDuration)
	fmt.Printf("   🐌 Slowest Step: %v\n", maxDuration)
	fmt.Printf("   🔧 Concurrency Used: %d\n", pe.config.MaxConcurrency)
}
