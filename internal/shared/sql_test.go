package shared

import "testing"

func TestCreateInPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{"zero count", 0, "?"},
		{"negative count", -1, "?"},
		{"single count", 1, "?"},
		{"three count", 3, "?,?,?"},
		{"five count", 5, "?,?,?,?,?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateInPlaceholders(tt.count)
			if result != tt.expected {
				t.Errorf("CreateInPlaceholders(%d) = %v, want %v", tt.count, result, tt.expected)
			}
		})
	}
}
