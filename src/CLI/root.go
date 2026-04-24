package cli

import (
	"dbtool/CLI/commands"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbtool",
	Short: "main",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return commands.InitConfig(cmd)
	},
}

func init() {
	rootCmd.AddCommand(commands.ConnectCmd())
	rootCmd.AddCommand(commands.BackupCmd())
	rootCmd.AddCommand(commands.RestoreCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
