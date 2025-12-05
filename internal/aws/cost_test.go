package aws

import (
	"testing"
)

func TestSortByAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    []CostResult
		expected []CostResult
	}{
		{
			name:     "empty slice",
			input:    []CostResult{},
			expected: []CostResult{},
		},
		{
			name: "single item",
			input: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
			},
		},
		{
			name: "already sorted",
			input: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
				{Service: "S3", Amount: 50.0, Unit: "USD"},
				{Service: "Lambda", Amount: 10.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
				{Service: "S3", Amount: 50.0, Unit: "USD"},
				{Service: "Lambda", Amount: 10.0, Unit: "USD"},
			},
		},
		{
			name: "reverse order",
			input: []CostResult{
				{Service: "Lambda", Amount: 10.0, Unit: "USD"},
				{Service: "S3", Amount: 50.0, Unit: "USD"},
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
				{Service: "S3", Amount: 50.0, Unit: "USD"},
				{Service: "Lambda", Amount: 10.0, Unit: "USD"},
			},
		},
		{
			name: "mixed order",
			input: []CostResult{
				{Service: "S3", Amount: 50.0, Unit: "USD"},
				{Service: "Lambda", Amount: 10.0, Unit: "USD"},
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
				{Service: "RDS", Amount: 75.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
				{Service: "RDS", Amount: 75.0, Unit: "USD"},
				{Service: "S3", Amount: 50.0, Unit: "USD"},
				{Service: "Lambda", Amount: 10.0, Unit: "USD"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortByAmount(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i].Service != tt.expected[i].Service {
					t.Errorf("index %d: service got %s, want %s", i, result[i].Service, tt.expected[i].Service)
				}
				if result[i].Amount != tt.expected[i].Amount {
					t.Errorf("index %d: amount got %f, want %f", i, result[i].Amount, tt.expected[i].Amount)
				}
			}
		})
	}
}
