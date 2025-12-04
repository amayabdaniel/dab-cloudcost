package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/amayabdaniel/dab-cloudcost/internal/aws"
	"github.com/spf13/cobra"
)

var (
	awsDays    int
	awsProfile string
	awsOutput  string
	awsTop     int
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Analyze AWS costs",
	Long:  `Fetch and analyze costs from AWS Cost Explorer API.`,
	RunE:  runAWS,
}

func runAWS(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	fmt.Printf("fetching aws costs for last %d days...\n\n", awsDays)

	client, err := aws.NewClient(ctx, awsProfile)
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	costs, err := client.GetCostsByService(ctx, awsDays)
	if err != nil {
		return fmt.Errorf("failed to get costs: %w", err)
	}

	if len(costs) == 0 {
		fmt.Println("no cost data found")
		return nil
	}

	if awsTop > 0 && awsTop < len(costs) {
		costs = costs[:awsTop]
	}

	switch awsOutput {
	case "json":
		return outputJSON(costs)
	case "csv":
		return outputCSV(costs)
	default:
		return outputTable(costs)
	}
}

func outputJSON(costs []aws.CostResult) error {
	var total float64
	for _, c := range costs {
		total += c.Amount
	}

	output := struct {
		Services []aws.CostResult `json:"services"`
		Total    float64          `json:"total"`
		Unit     string           `json:"unit"`
	}{
		Services: costs,
		Total:    total,
		Unit:     costs[0].Unit,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func outputCSV(costs []aws.CostResult) error {
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"service", "cost", "unit"})

	var total float64
	for _, c := range costs {
		w.Write([]string{c.Service, fmt.Sprintf("%.2f", c.Amount), c.Unit})
		total += c.Amount
	}

	w.Write([]string{"TOTAL", fmt.Sprintf("%.2f", total), costs[0].Unit})
	w.Flush()
	return w.Error()
}

func outputTable(costs []aws.CostResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tCOST\tUNIT")
	fmt.Fprintln(w, "-------\t----\t----")

	var total float64
	for _, c := range costs {
		fmt.Fprintf(w, "%s\t%.2f\t%s\n", c.Service, c.Amount, c.Unit)
		total += c.Amount
	}

	fmt.Fprintln(w, "-------\t----\t----")
	fmt.Fprintf(w, "TOTAL\t%.2f\t%s\n", total, costs[0].Unit)
	w.Flush()

	return nil
}

func init() {
	awsCmd.Flags().IntVarP(&awsDays, "days", "d", 30, "number of days to analyze")
	awsCmd.Flags().StringVarP(&awsProfile, "profile", "p", "default", "aws profile to use")
	awsCmd.Flags().StringVarP(&awsOutput, "output", "o", "table", "output format (table, json, csv)")
	awsCmd.Flags().IntVarP(&awsTop, "top", "t", 0, "show top N services (0 = all)")
	rootCmd.AddCommand(awsCmd)
}
