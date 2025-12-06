package gcp

import (
	"math"
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
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
			},
		},
		{
			name: "already sorted",
			input: []CostResult{
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
				{Service: "Cloud Storage", Amount: 50.0, Unit: "USD"},
				{Service: "Cloud Functions", Amount: 10.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
				{Service: "Cloud Storage", Amount: 50.0, Unit: "USD"},
				{Service: "Cloud Functions", Amount: 10.0, Unit: "USD"},
			},
		},
		{
			name: "reverse order",
			input: []CostResult{
				{Service: "Cloud Functions", Amount: 10.0, Unit: "USD"},
				{Service: "Cloud Storage", Amount: 50.0, Unit: "USD"},
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
				{Service: "Cloud Storage", Amount: 50.0, Unit: "USD"},
				{Service: "Cloud Functions", Amount: 10.0, Unit: "USD"},
			},
		},
		{
			name: "mixed order",
			input: []CostResult{
				{Service: "Cloud Storage", Amount: 50.0, Unit: "USD"},
				{Service: "Cloud Functions", Amount: 10.0, Unit: "USD"},
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
				{Service: "BigQuery", Amount: 75.0, Unit: "USD"},
			},
			expected: []CostResult{
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
				{Service: "BigQuery", Amount: 75.0, Unit: "USD"},
				{Service: "Cloud Storage", Amount: 50.0, Unit: "USD"},
				{Service: "Cloud Functions", Amount: 10.0, Unit: "USD"},
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

func TestTotalCost(t *testing.T) {
	tests := []struct {
		name     string
		input    []CostResult
		expected float64
	}{
		{
			name:     "empty slice",
			input:    []CostResult{},
			expected: 0.0,
		},
		{
			name: "single item",
			input: []CostResult{
				{Service: "Compute Engine", Amount: 100.50, Unit: "USD"},
			},
			expected: 100.50,
		},
		{
			name: "multiple items",
			input: []CostResult{
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
				{Service: "Cloud Storage", Amount: 50.25, Unit: "USD"},
				{Service: "Cloud Functions", Amount: 10.75, Unit: "USD"},
			},
			expected: 161.0,
		},
		{
			name: "with zero amounts",
			input: []CostResult{
				{Service: "Compute Engine", Amount: 100.0, Unit: "USD"},
				{Service: "Cloud Storage", Amount: 0.0, Unit: "USD"},
				{Service: "Cloud Functions", Amount: 50.0, Unit: "USD"},
			},
			expected: 150.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TotalCost(tt.input)
			if math.Abs(result-tt.expected) > 0.001 {
				t.Errorf("got %f, want %f", result, tt.expected)
			}
		})
	}
}
