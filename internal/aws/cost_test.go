package aws

import (
	"math"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
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
				{Service: "EC2", Amount: 100.50, Unit: "USD"},
			},
			expected: 100.50,
		},
		{
			name: "multiple items",
			input: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
				{Service: "S3", Amount: 50.25, Unit: "USD"},
				{Service: "Lambda", Amount: 10.75, Unit: "USD"},
			},
			expected: 161.0,
		},
		{
			name: "with zero amounts",
			input: []CostResult{
				{Service: "EC2", Amount: 100.0, Unit: "USD"},
				{Service: "S3", Amount: 0.0, Unit: "USD"},
				{Service: "Lambda", Amount: 50.0, Unit: "USD"},
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

func TestParseCostResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *costexplorer.GetCostAndUsageOutput
		expected []CostResult
	}{
		{
			name: "empty response",
			input: &costexplorer.GetCostAndUsageOutput{
				ResultsByTime: []types.ResultByTime{},
			},
			expected: nil,
		},
		{
			name: "single service",
			input: &costexplorer.GetCostAndUsageOutput{
				ResultsByTime: []types.ResultByTime{
					{
						Groups: []types.Group{
							{
								Keys: []string{"Amazon EC2"},
								Metrics: map[string]types.MetricValue{
									"UnblendedCost": {
										Amount: aws.String("150.75"),
										Unit:   aws.String("USD"),
									},
								},
							},
						},
					},
				},
			},
			expected: []CostResult{
				{Service: "Amazon EC2", Amount: 150.75, Unit: "USD"},
			},
		},
		{
			name: "multiple services sorted",
			input: &costexplorer.GetCostAndUsageOutput{
				ResultsByTime: []types.ResultByTime{
					{
						Groups: []types.Group{
							{
								Keys: []string{"Amazon S3"},
								Metrics: map[string]types.MetricValue{
									"UnblendedCost": {Amount: aws.String("50.00"), Unit: aws.String("USD")},
								},
							},
							{
								Keys: []string{"Amazon EC2"},
								Metrics: map[string]types.MetricValue{
									"UnblendedCost": {Amount: aws.String("200.00"), Unit: aws.String("USD")},
								},
							},
							{
								Keys: []string{"AWS Lambda"},
								Metrics: map[string]types.MetricValue{
									"UnblendedCost": {Amount: aws.String("25.00"), Unit: aws.String("USD")},
								},
							},
						},
					},
				},
			},
			expected: []CostResult{
				{Service: "Amazon EC2", Amount: 200.00, Unit: "USD"},
				{Service: "Amazon S3", Amount: 50.00, Unit: "USD"},
				{Service: "AWS Lambda", Amount: 25.00, Unit: "USD"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCostResponse(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i].Service != tt.expected[i].Service {
					t.Errorf("index %d: service got %s, want %s", i, result[i].Service, tt.expected[i].Service)
				}
				if math.Abs(result[i].Amount-tt.expected[i].Amount) > 0.001 {
					t.Errorf("index %d: amount got %f, want %f", i, result[i].Amount, tt.expected[i].Amount)
				}
				if result[i].Unit != tt.expected[i].Unit {
					t.Errorf("index %d: unit got %s, want %s", i, result[i].Unit, tt.expected[i].Unit)
				}
			}
		})
	}
}
