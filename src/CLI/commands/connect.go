package commands

import (
	"dbtool/DBinterface"
	"dbtool/DBinterface/DBdrivers"
	"fmt"

	"github.com/spf13/cobra"
)

func ConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "connect",
		Short:        "command for connect to db",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			conf := DBinterface.Config{}
			if err := InitCmd(cmd, &conf); err != nil {
				return err
			}

			fmt.Printf("Connecting to: %s@%s:%d/%s\n", conf.User, conf.Host, conf.Port, conf.DBName)
			switch conf.DBtype {
			case "postgres":
				{
					if RuntimeErr := DBdrivers.ConnectToPostgres(conf); RuntimeErr != nil {
						return RuntimeErr
					}
				}
			case "mysql":
				{
					if RuntimeErr := DBdrivers.ConnectToMySQL(conf); RuntimeErr != nil {
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

	cmd.Flags().StringP("host", "H", "", "Database host")
	cmd.Flags().IntP("port", "P", 0, "Database port")
	cmd.Flags().StringP("database", "", "", "Type DB")
	cmd.Flags().StringP("dbname", "N", "", "Database name")
	cmd.Flags().StringP("username", "U", "", "Username")
	cmd.Flags().StringP("password", "p", "", "Password")

	return cmd
}
