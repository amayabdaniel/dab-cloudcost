package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/amayabdaniel/dab-cloudcost/internal/aws"
	"github.com/spf13/cobra"
)

var (
	awsDays    int
	awsProfile string
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
	rootCmd.AddCommand(awsCmd)
}
