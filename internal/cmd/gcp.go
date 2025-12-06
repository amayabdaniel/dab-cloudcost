package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/amayabdaniel/dab-cloudcost/internal/gcp"
	"github.com/spf13/cobra"
)

var (
	gcpDays    int
	gcpProject string
	gcpOutput  string
	gcpTop     int
	gcpTable   string
)

var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Analyze GCP costs",
	Long:  `Fetch and analyze costs from GCP BigQuery billing export.`,
	RunE:  runGCP,
}

func runGCP(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	fmt.Printf("fetching gcp costs for project '%s' (last %d days)...\n\n", gcpProject, gcpDays)

	client, err := gcp.NewClient(ctx, gcpProject, gcpTable)
	if err != nil {
		return fmt.Errorf("failed to create gcp client: %w", err)
	}
	defer client.Close()

	costs, err := client.GetCostsByService(ctx, gcpDays)
	if err != nil {
		return fmt.Errorf("failed to get costs: %w", err)
	}

	if len(costs) == 0 {
		fmt.Println("no cost data found")
		return nil
	}

	if gcpTop > 0 && gcpTop < len(costs) {
		costs = costs[:gcpTop]
	}

	switch gcpOutput {
	case "json":
		return gcpOutputJSON(costs)
	case "csv":
		return gcpOutputCSV(costs)
	default:
		return gcpOutputTable(costs)
	}
}

func gcpOutputJSON(costs []gcp.CostResult) error {
	var total float64
	for _, c := range costs {
		total += c.Amount
	}

	output := struct {
		Services []gcp.CostResult `json:"services"`
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

func gcpOutputCSV(costs []gcp.CostResult) error {
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

func gcpOutputTable(costs []gcp.CostResult) error {
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
	gcpCmd.Flags().IntVarP(&gcpDays, "days", "d", 30, "number of days to analyze")
	gcpCmd.Flags().StringVarP(&gcpProject, "project", "p", "", "gcp project id (required)")
	gcpCmd.Flags().StringVarP(&gcpOutput, "output", "o", "table", "output format (table, json, csv)")
	gcpCmd.Flags().IntVarP(&gcpTop, "top", "t", 0, "show top N services (0 = all)")
	gcpCmd.Flags().StringVar(&gcpTable, "billing-table", "", "bigquery billing export table (e.g. project.dataset.table)")
	gcpCmd.MarkFlagRequired("project")
	gcpCmd.MarkFlagRequired("billing-table")
	rootCmd.AddCommand(gcpCmd)
}
