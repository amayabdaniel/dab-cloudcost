package gcp

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type CostResult struct {
	Service string  `json:"service"`
	Amount  float64 `json:"amount"`
	Unit    string  `json:"unit"`
}

// BigQueryAPI interface for testing
type BigQueryAPI interface {
	Query(q string) *bigquery.Query
}

type Client struct {
	bq           *bigquery.Client
	projectID    string
	billingTable string
}

func NewClient(ctx context.Context, projectID, billingTable string) (*Client, error) {
	bq, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create bigquery client: %w", err)
	}

	return &Client{
		bq:           bq,
		projectID:    projectID,
		billingTable: billingTable,
	}, nil
}

func (c *Client) Close() error {
	return c.bq.Close()
}

func (c *Client) GetCostsByService(ctx context.Context, days int) ([]CostResult, error) {
	query := fmt.Sprintf(`
		SELECT
			service.description AS service,
			SUM(cost) AS amount,
			currency AS unit
		FROM %s
		WHERE DATE(_PARTITIONTIME) >= DATE_SUB(CURRENT_DATE(), INTERVAL %d DAY)
			AND cost > 0
		GROUP BY service.description, currency
		ORDER BY amount DESC
	`, c.billingTable, days)

	q := c.bq.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	var results []CostResult
	for {
		var row struct {
			Service string  `bigquery:"service"`
			Amount  float64 `bigquery:"amount"`
			Unit    string  `bigquery:"unit"`
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}
		results = append(results, CostResult{
			Service: row.Service,
			Amount:  row.Amount,
			Unit:    row.Unit,
		})
	}

	return SortByAmount(results), nil
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
