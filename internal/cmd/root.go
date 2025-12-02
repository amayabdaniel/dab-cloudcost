package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "dab-cloudcost",
	Short: "Cloud cost analyzer for AWS and GCP",
	Long: `dab-cloudcost analyzes your cloud costs across AWS and GCP,
providing insights and recommendations for cost optimization.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dab-cloudcost.yaml)")
}
