package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Performance benchmarks for CLI commands

// Helper function to create a fresh root command for tests
func getPerfTestRootCommand() *cobra.Command {
	factory := NewCommandFactory()
	return factory.NewRootCommand()
}

// BenchmarkInitCommand benchmarks the init command
func BenchmarkInitCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tempDir := b.TempDir()
		originalDir, _ := os.Getwd()
		os.Chdir(tempDir)

		rootCmd := getPerfTestRootCommand()

		b.StartTimer()
		executeCommand(rootCmd, "init")
		b.StopTimer()

		os.Chdir(originalDir)
	}
}

// BenchmarkListCommand benchmarks the list command with various resource counts
func BenchmarkListCommand(b *testing.B) {
	benchmarks := []struct {
		name          string
		resourceCount int
	}{
		{"10_resources", 10},
		{"100_resources", 100},
		{"1000_resources", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Setup test environment
			homeDir := b.TempDir()
			b.Setenv("HOME", homeDir)
			ddxHome := filepath.Join(homeDir, ".ddx")

			// Create resources
			templatesDir := filepath.Join(ddxHome, "templates")
			for i := 0; i < bm.resourceCount; i++ {
				dir := filepath.Join(templatesDir, fmt.Sprintf("template%d", i))
				os.MkdirAll(dir, 0755)
				os.WriteFile(filepath.Join(dir, "README.md"), []byte("Test"), 0644)
			}

			rootCmd := getPerfTestRootCommand()
			// Commands already registered in factory

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				executeCommand(rootCmd, "list")
			}
		})
	}
}

// BenchmarkApplyTemplate benchmarks template application
func BenchmarkApplyTemplate(b *testing.B) {
	benchmarks := []struct {
		name      string
		fileCount int
		fileSize  int
	}{
		{"small_template", 5, 1024},     // 5 files, 1KB each
		{"medium_template", 50, 10240},  // 50 files, 10KB each
		{"large_template", 200, 102400}, // 200 files, 100KB each
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Setup template
			homeDir := b.TempDir()
			b.Setenv("HOME", homeDir)

			templateDir := filepath.Join(homeDir, ".ddx", "templates", "perf-test")
			os.MkdirAll(templateDir, 0755)

			// Create template files
			for i := 0; i < bm.fileCount; i++ {
				content := strings.Repeat("x", bm.fileSize)
				filePath := filepath.Join(templateDir, fmt.Sprintf("file%d.txt", i))
				os.WriteFile(filePath, []byte(content), 0644)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				workDir := b.TempDir()
				os.Chdir(workDir)

				// Create config
				config := `version: "1.0"`
				os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644)

				rootCmd := getPerfTestRootCommand()
				// Commands already registered in factory

				b.StartTimer()
				executeCommand(rootCmd, "apply", "templates/perf-test")
			}
		})
	}
}

// TestPerformance_MemoryUsage tests memory consumption
func TestPerformance_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory test in short mode")
	}

	tests := []struct {
		name        string
		operation   func()
		maxMemoryMB uint64
	}{
		{
			name: "list_command_memory",
			operation: func() {
				rootCmd := getPerfTestRootCommand()
				// Commands already registered in factory
				executeCommand(rootCmd, "list")
			},
			maxMemoryMB: 50, // Max 50MB for list operation
		},
		{
			name: "config_command_memory",
			operation: func() {
				workDir := t.TempDir()
				os.Chdir(workDir)

				config := `version: "1.0"`
				os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644)

				rootCmd := getPerfTestRootCommand()
				// Commands already registered in factory
				executeCommand(rootCmd, "config")
			},
			maxMemoryMB: 30, // Max 30MB for config operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get initial memory stats
			var m1 runtime.MemStats
			runtime.ReadMemStats(&m1)

			// Run operation
			tt.operation()

			// Force GC and get final memory stats
			runtime.GC()
			var m2 runtime.MemStats
			runtime.ReadMemStats(&m2)

			// Calculate memory used (handle cases where GC freed memory)
			var memUsedMB uint64
			if m2.Alloc > m1.Alloc {
				memUsedMB = (m2.Alloc - m1.Alloc) / 1024 / 1024
			} else {
				// If GC freed memory, consider it as 0 additional memory used
				memUsedMB = 0
			}

			assert.LessOrEqual(t, memUsedMB, tt.maxMemoryMB,
				"Memory usage exceeds limit: %dMB > %dMB", memUsedMB, tt.maxMemoryMB)
		})
	}
}

// TestPerformance_ResponseTime tests command response times
func TestPerformance_ResponseTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping response time test in short mode")
	}

	tests := []struct {
		name        string
		args        []string
		setup       func() string
		maxDuration time.Duration
	}{
		{
			name:        "help_response_time",
			args:        []string{"--help"},
			maxDuration: 100 * time.Millisecond,
		},
		{
			name:        "version_response_time",
			args:        []string{"--version"},
			maxDuration: 50 * time.Millisecond,
		},
		{
			name: "list_response_time",
			args: []string{"list"},
			setup: func() string {
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				return homeDir
			},
			maxDuration: 200 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			rootCmd := getPerfTestRootCommand()
			// Commands already registered
			// Commands already registered in factory
			// Commands already registered
			// Commands already registered

			start := time.Now()
			executeCommand(rootCmd, tt.args...)
			duration := time.Since(start)

			assert.LessOrEqual(t, duration, tt.maxDuration,
				"Response time exceeds limit: %v > %v", duration, tt.maxDuration)
		})
	}
}

// TestPerformance_ConcurrentOperations tests concurrent command execution
func TestPerformance_ConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent operations test in short mode")
	}

	const numWorkers = 10
	const numOperations = 100

	// Setup environment
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	// Create some templates
	templatesDir := filepath.Join(homeDir, ".ddx", "templates")
	for i := 0; i < 5; i++ {
		dir := filepath.Join(templatesDir, fmt.Sprintf("template%d", i))
		os.MkdirAll(dir, 0755)
	}

	var wg sync.WaitGroup
	errors := make(chan error, numWorkers*numOperations)
	start := time.Now()

	// Launch concurrent workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				rootCmd := getPerfTestRootCommand()
				// Commands already registered in factory

				_, err := executeCommand(rootCmd, "list")
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)
	duration := time.Since(start)

	// Check for errors
	errorCount := 0
	for err := range errors {
		errorCount++
		t.Logf("Concurrent operation error: %v", err)
	}

	// Calculate throughput
	totalOps := numWorkers * numOperations
	throughput := float64(totalOps) / duration.Seconds()

	t.Logf("Completed %d operations in %v (%.2f ops/sec)", totalOps, duration, throughput)

	// Assertions
	assert.Less(t, errorCount, totalOps/10, "Error rate should be < 10%")
	assert.Greater(t, throughput, 100.0, "Throughput should be > 100 ops/sec")
}

// TestPerformance_LargeFileHandling tests handling of large files
func TestPerformance_LargeFileHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file test in short mode")
	}

	sizes := []struct {
		name     string
		fileSize int
		maxTime  time.Duration
	}{
		{"1MB", 1024 * 1024, 500 * time.Millisecond},
		{"10MB", 10 * 1024 * 1024, 2 * time.Second},
		{"50MB", 50 * 1024 * 1024, 5 * time.Second},
	}

	for _, size := range sizes {
		t.Run(size.name, func(t *testing.T) {
			// Create large template file
			homeDir := t.TempDir()
			t.Setenv("HOME", homeDir)

			templateDir := filepath.Join(homeDir, ".ddx", "templates", "large")
			os.MkdirAll(templateDir, 0755)

			largeFile := filepath.Join(templateDir, "large.txt")
			content := make([]byte, size.fileSize)
			for i := range content {
				content[i] = byte('A' + (i % 26))
			}
			os.WriteFile(largeFile, content, 0644)

			// Setup work directory
			workDir := t.TempDir()
			os.Chdir(workDir)

			config := `version: "1.0"`
			os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644)

			rootCmd := getPerfTestRootCommand()
			// Commands already registered

			// Measure time to apply large template
			start := time.Now()
			executeCommand(rootCmd, "apply", "templates/large")
			duration := time.Since(start)

			assert.LessOrEqual(t, duration, size.maxTime,
				"Processing %s file took too long: %v > %v", size.name, duration, size.maxTime)
		})
	}
}

// TestPerformance_StartupTime tests CLI startup overhead
func TestPerformance_StartupTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping startup time test in short mode")
	}

	const maxStartupTime = 50 * time.Millisecond
	const iterations = 10

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()

		rootCmd := getPerfTestRootCommand()
		// Commands already registered
		// Commands already registered
		// Commands already registered

		// Just initialize, don't execute
		rootCmd.InitDefaultHelpCmd()

		duration := time.Since(start)
		totalTime += duration
	}

	avgTime := totalTime / iterations

	assert.LessOrEqual(t, avgTime, maxStartupTime,
		"Average startup time exceeds limit: %v > %v", avgTime, maxStartupTime)
}

// BenchmarkVariableSubstitution benchmarks template variable replacement
func BenchmarkVariableSubstitution(b *testing.B) {
	benchmarks := []struct {
		name          string
		variableCount int
		templateSize  int
	}{
		{"10_vars_small", 10, 1024},
		{"50_vars_medium", 50, 10240},
		{"100_vars_large", 100, 102400},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Create template with variables
			var template strings.Builder
			for i := 0; i < bm.templateSize; i += 100 {
				template.WriteString(fmt.Sprintf("Line with {{var%d}} variable\n", i%bm.variableCount))
			}

			// Create variables map
			variables := make(map[string]string)
			for i := 0; i < bm.variableCount; i++ {
				variables[fmt.Sprintf("var%d", i)] = fmt.Sprintf("value%d", i)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Benchmark variable substitution
				// This would call the actual template processing function
				_ = processTemplate(template.String(), variables)
			}
		})
	}
}

// Helper function to simulate template processing
func processTemplate(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}

// TestPerformance_ResourceCleanup tests proper resource cleanup
func TestPerformance_ResourceCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping resource cleanup test in short mode")
	}

	// Get initial goroutine count
	initialGoroutines := runtime.NumGoroutine()

	// Run multiple operations
	for i := 0; i < 100; i++ {
		rootCmd := getPerfTestRootCommand()
		// Commands already registered
		executeCommand(rootCmd, "list")
	}

	// Allow time for cleanup
	time.Sleep(100 * time.Millisecond)
	runtime.GC()

	// Check goroutine count
	finalGoroutines := runtime.NumGoroutine()

	// Allow for some variance, but detect leaks
	assert.LessOrEqual(t, finalGoroutines, initialGoroutines+5,
		"Possible goroutine leak: %d -> %d", initialGoroutines, finalGoroutines)
}
