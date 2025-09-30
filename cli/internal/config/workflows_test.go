package config

import (
	"testing"
)

// TestWorkflowsConfig_ApplyDefaults tests default value assignment
func TestWorkflowsConfig_ApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		initial  WorkflowsConfig
		expected WorkflowsConfig
	}{
		{
			name:    "empty config gets defaults",
			initial: WorkflowsConfig{},
			expected: WorkflowsConfig{
				Active:   []string{},
				SafeWord: "NODDX",
			},
		},
		{
			name: "existing safe word preserved",
			initial: WorkflowsConfig{
				SafeWord: "CUSTOM",
			},
			expected: WorkflowsConfig{
				Active:   []string{},
				SafeWord: "CUSTOM",
			},
		},
		{
			name: "existing active preserved",
			initial: WorkflowsConfig{
				Active: []string{"helix"},
			},
			expected: WorkflowsConfig{
				Active:   []string{"helix"},
				SafeWord: "NODDX",
			},
		},
		{
			name: "nil active becomes empty slice",
			initial: WorkflowsConfig{
				Active: nil,
			},
			expected: WorkflowsConfig{
				Active:   []string{},
				SafeWord: "NODDX",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.initial
			cfg.ApplyDefaults()

			if cfg.SafeWord != tt.expected.SafeWord {
				t.Errorf("SafeWord = %v, want %v", cfg.SafeWord, tt.expected.SafeWord)
			}

			if len(cfg.Active) != len(tt.expected.Active) {
				t.Errorf("Active length = %v, want %v", len(cfg.Active), len(tt.expected.Active))
			}
		})
	}
}

// TestWorkflowsConfig_Validate tests validation logic
func TestWorkflowsConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  WorkflowsConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: WorkflowsConfig{
				Active:   []string{"helix"},
				SafeWord: "NODDX",
			},
			wantErr: false,
		},
		{
			name: "empty safe word",
			config: WorkflowsConfig{
				Active:   []string{"helix"},
				SafeWord: "",
			},
			wantErr: true,
			errMsg:  "safe_word cannot be empty",
		},
		{
			name: "safe word with spaces",
			config: WorkflowsConfig{
				Active:   []string{"helix"},
				SafeWord: "NO DDX",
			},
			wantErr: true,
			errMsg:  "safe_word cannot contain spaces",
		},
		{
			name: "duplicate workflows",
			config: WorkflowsConfig{
				Active:   []string{"helix", "kanban", "helix"},
				SafeWord: "NODDX",
			},
			wantErr: true,
			errMsg:  "duplicate workflow",
		},
		{
			name: "empty active list valid",
			config: WorkflowsConfig{
				Active:   []string{},
				SafeWord: "NODDX",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestNewConfig_WithWorkflows tests that NewConfig includes workflows field
func TestNewConfig_WithWorkflows(t *testing.T) {
	cfg := DefaultNewConfig()
	cfg.ApplyDefaults()

	if cfg.Workflows.SafeWord != "NODDX" {
		t.Errorf("default SafeWord = %v, want NODDX", cfg.Workflows.SafeWord)
	}

	if cfg.Workflows.Active == nil {
		t.Error("Active should not be nil after ApplyDefaults")
	}
}

// TestConfig_LoadWithWorkflows tests loading config with workflows section
func TestConfig_LoadWithWorkflows(t *testing.T) {
	// This will be an integration test once loader is implemented
	t.Skip("requires config loader implementation")
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
