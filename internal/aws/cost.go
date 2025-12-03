package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type CostResult struct {
	Service string
	Amount  string
	Unit    string
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
				results = append(results, CostResult{
					Service: group.Keys[0],
					Amount:  aws.ToString(cost.Amount),
					Unit:    aws.ToString(cost.Unit),
				})
			}
		}
	}

	return results, nil
}
