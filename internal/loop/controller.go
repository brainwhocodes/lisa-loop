package loop

import (
	stdcontext "context"
	"fmt"
	"strings"
	"time"

	"github.com/brainwhocodes/ralph-codex/internal/circuit"
	"github.com/brainwhocodes/ralph-codex/internal/codex"
)

// Config holds configuration for the loop
type Config struct {
	Backend      string
	ProjectPath  string
	PromptPath   string
	MaxCalls     int
	Timeout      int
	Verbose      bool
	ResetCircuit bool
}

// Controller manages the main Ralph loop
type Controller struct {
	config      ControllerConfig
	rateLimiter *RateLimiter
	breaker     *circuit.Breaker
	codexRunner *codex.Runner
	loopNum     int
	lastOutput  string
	shouldStop  bool
}

// ControllerConfig holds configuration for the loop controller
type ControllerConfig struct {
	MaxLoops      int
	MaxDuration   time.Duration
	CheckInterval time.Duration
}

// NewController creates a new loop controller
func NewController(config Config, rateLimiter *RateLimiter, breaker *circuit.Breaker) *Controller {
	codexConfig := codex.Config{
		Backend:      config.Backend,
		ProjectPath:  config.ProjectPath,
		PromptPath:   config.PromptPath,
		MaxCalls:     config.MaxCalls,
		Timeout:      config.Timeout,
		Verbose:      config.Verbose,
		ResetCircuit: config.ResetCircuit,
	}
	codexRunner := codex.NewRunner(codexConfig)

	return &Controller{
		config: ControllerConfig{
			MaxLoops:      config.MaxCalls,
			MaxDuration:   time.Duration(config.Timeout) * time.Second,
			CheckInterval: 5 * time.Second,
		},
		rateLimiter: rateLimiter,
		breaker:     breaker,
		codexRunner: codexRunner,
		loopNum:     0,
		lastOutput:  "",
		shouldStop:  false,
	}
}

// Run executes the main loop
func (c *Controller) Run(ctx stdcontext.Context) error {
	fmt.Printf("\n🚀 Starting Ralph Codex loop (max %d calls)...\n", c.config.MaxLoops)

	for {
		if c.shouldStop {
			fmt.Println("\n✅ Loop stopped")
			return nil
		}

		select {
		case <-ctx.Done():
			fmt.Println("\n🛑 Loop cancelled")
			return ctx.Err()
		default:
			// Execute one iteration
			err := c.ExecuteLoop(ctx)

			if err != nil {
				fmt.Printf("\n❌ Loop iteration error: %v\n", err)
				return err
			}

			// Check if we should stop
			if c.ShouldContinue() {
				fmt.Printf("\n✅ Ralph Codex loop complete after %d iterations\n", c.loopNum)
				return nil
			}

			c.loopNum++
		}
	}
}

// ExecuteLoop executes a single loop iteration
func (c *Controller) ExecuteLoop(ctx stdcontext.Context) error {
	// Check rate limit
	if !c.rateLimiter.CanMakeCall() {
		fmt.Printf("\n⏱️  Rate limit reached. Calls remaining: %d\n", c.rateLimiter.CallsRemaining())
		return c.rateLimiter.WaitForReset(ctx)
	}

	// Check circuit breaker
	if c.breaker.ShouldHalt() {
		return fmt.Errorf("circuit breaker is OPEN, halting execution")
	}

	// Load prompt and fix plan
	prompt, err := GetPrompt()
	if err != nil {
		return fmt.Errorf("failed to load prompt: %w", err)
	}

	tasks, err := LoadFixPlan()
	if err != nil {
		return fmt.Errorf("failed to load fix plan: %w", err)
	}

	// Build context
	circuitState := c.breaker.GetState().String()
	remainingTasks := []string{}
	for _, task := range tasks {
		if !strings.HasPrefix(task, "[x]") {
			remainingTasks = append(remainingTasks, task)
		}
	}

	loopContext, err := BuildContext("", c.loopNum+1, remainingTasks, circuitState, c.lastOutput)
	if err != nil {
		return fmt.Errorf("failed to build context: %w", err)
	}

	promptWithContext := InjectContext(prompt, loopContext)

	// Execute Codex
	fmt.Printf("\n🔄 Loop %d: Executing Codex...\n", c.loopNum+1)
	output, _, err := c.codexRunner.Run(promptWithContext)

	if err != nil {
		c.lastOutput = fmt.Sprintf("Error: %v", err)
		c.rateLimiter.RecordCall()

		// Record error in circuit breaker
		c.breaker.RecordError(err.Error())
		return err
	}

	c.lastOutput = fmt.Sprintf("Success: %s", output[:min(200, len(output))])

	// Analyze output for exit conditions
	// TODO: This will be implemented in response analysis package

	// Record result in circuit breaker
	filesChanged := 0
	if strings.Contains(output, "Modified") || strings.Contains(output, "Created") {
		filesChanged = 1
	}

	hasErrors := strings.Contains(output, "Error") || strings.Contains(output, "failed")

	err = c.breaker.RecordResult(c.loopNum, filesChanged, hasErrors)
	if err != nil {
		return err
	}

	return nil
}

// ShouldContinue checks if the loop should continue
func (c *Controller) ShouldContinue() bool {
	tasks, err := LoadFixPlan()
	if err != nil {
		return false
	}

	// Check if all tasks are complete
	allComplete := true
	for _, task := range tasks {
		if !strings.HasPrefix(task, "[x]") {
			allComplete = false
			break
		}
	}

	if allComplete {
		c.shouldStop = true
		return true
	}

	// Check circuit breaker
	if c.breaker.ShouldHalt() {
		c.shouldStop = true
		return true
	}

	// Check rate limit
	if !c.rateLimiter.CanMakeCall() {
		c.shouldStop = true
		return true
	}

	// Check max loops
	if c.loopNum >= c.config.MaxLoops {
		c.shouldStop = true
		return true
	}

	return false
}

// CheckExitConditions analyzes output for completion signals
func (c *Controller) CheckExitConditions(output string) bool {
	// TODO: Will be implemented with response analysis package
	return false
}

// HandleCircuitBreakerOpen handles circuit breaker being open
func (c *Controller) HandleCircuitBreakerOpen() error {
	return fmt.Errorf("circuit breaker is %s", c.breaker.GetState())
}

// HandleRateLimitExceeded handles rate limit being exceeded
func (c *Controller) HandleRateLimitExceeded() error {
	return c.rateLimiter.WaitForReset(stdcontext.Background())
}

// UpdateProgress updates loop progress for display
func (c *Controller) UpdateProgress(loopNum int, status string) {
	// This would update TUI model in the full implementation
	// For now, just print
	fmt.Printf("📊 Progress: Loop %d - %s\n", loopNum, status)
}

// GracefulExit performs cleanup before exiting
func (c *Controller) GracefulExit() error {
	fmt.Println("\n🧹 Performing graceful exit...")

	c.shouldStop = true

	// Reset circuit breaker
	c.breaker.Reset()

	// Reset session
	if err := codex.NewSession(); err != nil {
		return fmt.Errorf("failed to reset session: %w", err)
	}

	fmt.Println("✅ Graceful exit complete")
	return nil
}

// GetStats returns controller statistics
func (c *Controller) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"loop_num":        c.loopNum,
		"should_stop":     c.shouldStop,
		"rate_limiter":    c.rateLimiter.GetStats(),
		"circuit_breaker": c.breaker.GetStats(),
		"last_output":     c.lastOutput,
	}
}
