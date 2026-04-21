package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "db-backup",
	Short: "main",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return InitConfig(cmd)
	},
}

func init() {
	rootCmd.AddCommand()
}
