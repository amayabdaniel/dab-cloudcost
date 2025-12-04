package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gcpDays    int
	gcpProject string
	gcpOutput  string
	gcpTop     int
)

var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Analyze GCP costs",
	Long:  `Fetch and analyze costs from GCP Cloud Billing API.`,
	RunE:  runGCP,
}

func runGCP(cmd *cobra.Command, args []string) error {
	fmt.Printf("fetching gcp costs for project '%s' (last %d days)...\n\n", gcpProject, gcpDays)
	fmt.Println("gcp billing integration coming soon")
	return nil
}

func init() {
	gcpCmd.Flags().IntVarP(&gcpDays, "days", "d", 30, "number of days to analyze")
	gcpCmd.Flags().StringVarP(&gcpProject, "project", "p", "", "gcp project id (required)")
	gcpCmd.Flags().StringVarP(&gcpOutput, "output", "o", "table", "output format (table, json, csv)")
	gcpCmd.Flags().IntVarP(&gcpTop, "top", "t", 0, "show top N services (0 = all)")
	gcpCmd.MarkFlagRequired("project")
	rootCmd.AddCommand(gcpCmd)
}
