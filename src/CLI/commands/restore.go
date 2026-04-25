package commands

import (
	"dbtool/DBinterface"
	"dbtool/DBinterface/DBdrivers"
	"fmt"

	"github.com/spf13/cobra"
)

func RestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "restore",
		Short:        "command to restore DB",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			conf := DBinterface.Config{}
			if err := InitCmd(cmd, &conf); err != nil {
				return err
			}
			conf.FilePath = GlobalViper.GetString("path")

			switch conf.DBtype {
			case "postgres":
				{
					if RuntimeErr := DBdrivers.RestorePostgres(conf); RuntimeErr != nil {
						return RuntimeErr
					}
				}
			case "mysql":
				{
					if RuntimeErr := DBdrivers.RestoreMySQL(conf); RuntimeErr != nil {
						return RuntimeErr
					}
				}
			default:
				{
					return fmt.Errorf("Unknown DB name!")
				}
			}
			return nil
		},
	}
	cmd.Flags().StringP("host", "H", "", "Host")
	cmd.Flags().IntP("port", "P", 0, "Port")
	cmd.Flags().StringP("database", "", "postgres", "Type DB")
	cmd.Flags().StringP("dbname", "N", "", "Name db to connect")
	cmd.Flags().StringP("username", "U", "", "Username")
	cmd.Flags().StringP("password", "", "", "Password")
	cmd.Flags().StringP("path", "", ".", "The path where file will be checked")

	return cmd
}
