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

type Client struct {
	ce *costexplorer.Client
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

	sort.Slice(results, func(i, j int) bool {
		return results[i].Amount > results[j].Amount
	})

	return results, nil
}
