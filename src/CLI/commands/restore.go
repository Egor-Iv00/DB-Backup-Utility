package commands

import (
	"github.com/spf13/cobra"
)

func RestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "command to restore DB",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	cmd.Flags().StringP("host", "H", "", "Host")
	cmd.Flags().IntP("port", "P", 0, "Port")
	cmd.Flags().StringP("database", "", "postgres", "Type DB")
	cmd.Flags().StringP("dbname", "N", "", "Name db to connect")
	cmd.Flags().StringP("username", "U", "", "Username")
	cmd.Flags().StringP("password", "", "", "Password")
	return cmd
}
