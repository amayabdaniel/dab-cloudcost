package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("fetching aws costs for last %d days using profile '%s'...\n", awsDays, awsProfile)
		// TODO: implement cost fetching
		return nil
	},
}

func init() {
	awsCmd.Flags().IntVarP(&awsDays, "days", "d", 30, "number of days to analyze")
	awsCmd.Flags().StringVarP(&awsProfile, "profile", "p", "default", "aws profile to use")
	rootCmd.AddCommand(awsCmd)
}
