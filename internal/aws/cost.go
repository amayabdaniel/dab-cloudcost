package aws

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type CostResult struct {
	Service string  `json:"service"`
	Amount  float64 `json:"amount"`
	Unit    string  `json:"unit"`
}

// CostExplorerAPI interface for testing
type CostExplorerAPI interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

type Client struct {
	ce CostExplorerAPI
}

func NewClient(ctx context.Context, profile string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		ce: costexplorer.NewFromConfig(cfg),
	}, nil
}

// NewClientWithAPI creates a client with a custom API (for testing)
func NewClientWithAPI(api CostExplorerAPI) *Client {
	return &Client{ce: api}
}

func (c *Client) GetCostsByService(ctx context.Context, days int) ([]CostResult, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -days)

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(start.Format("2006-01-02")),
			End:   aws.String(end.Format("2006-01-02")),
		},
		Granularity: types.GranularityMonthly,
		Metrics:     []string{"UnblendedCost"},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
		},
	}

	output, err := c.ce.GetCostAndUsage(ctx, input)
	if err != nil {
		return nil, err
	}

	return ParseCostResponse(output), nil
}

// ParseCostResponse parses AWS cost response into CostResults
func ParseCostResponse(output *costexplorer.GetCostAndUsageOutput) []CostResult {
	var results []CostResult
	for _, result := range output.ResultsByTime {
		for _, group := range result.Groups {
			if len(group.Keys) > 0 {
				cost := group.Metrics["UnblendedCost"]
				amount, _ := strconv.ParseFloat(aws.ToString(cost.Amount), 64)
				results = append(results, CostResult{
					Service: group.Keys[0],
					Amount:  amount,
					Unit:    aws.ToString(cost.Unit),
				})
			}
		}
	}
	return SortByAmount(results)
}

// SortByAmount sorts results by amount descending
func SortByAmount(results []CostResult) []CostResult {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Amount > results[j].Amount
	})
	return results
}

// TotalCost calculates total cost from results
func TotalCost(results []CostResult) float64 {
	var total float64
	for _, r := range results {
		total += r.Amount
	}
	return total
}
