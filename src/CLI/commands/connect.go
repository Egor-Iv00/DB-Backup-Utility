package commands

import (
	"dbtool/DBdrivers"

	"github.com/spf13/cobra"
)

func ConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "command for connect to db",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := DBdrivers.Config{}
			//typeDB:= cli.GlobalViper.GetString("database")
			conf.DBName = GlobalViper.GetString("dbname")
			conf.Host = GlobalViper.GetString("host")
			conf.Port = GlobalViper.GetInt("port")
			conf.User = GlobalViper.GetString("username")
			conf.Password = GlobalViper.GetString("password")

			if RuntimeErr := DBdrivers.ConnectToPostgres(conf); RuntimeErr != nil {
				return RuntimeErr
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
	return cmd
}
