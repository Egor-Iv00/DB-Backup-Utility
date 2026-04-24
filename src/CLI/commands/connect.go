package commands

import (
	"dbtool/DBdrivers"
	"fmt"

	"github.com/spf13/cobra"
)

func ConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "connect",
		Short:        "command for connect to db",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := InitConfig(cmd); err != nil {
				return fmt.Errorf("init config error: %w", err)
			}

			conf := DBdrivers.Config{}
			conf.Host = GlobalViper.GetString("host")
			conf.Port = GlobalViper.GetInt("port")
			conf.DBName = GlobalViper.GetString("dbname")
			conf.User = GlobalViper.GetString("username")
			conf.Password = GlobalViper.GetString("password")

			if conf.Host == "" {
				return fmt.Errorf("host is empty. Set it in config.json or use --host flag")
			}
			if conf.User == "" {
				return fmt.Errorf("username is empty. Set it in config.json or use --username flag")
			}
			if conf.Port == 0 {
				conf.Port = 5432
			}

			fmt.Printf("Connecting to: %s@%s:%d/%s\n", conf.User, conf.Host, conf.Port, conf.DBName)
			if RuntimeErr := DBdrivers.ConnectToPostgres(conf); RuntimeErr != nil {
				return RuntimeErr
			}
			return nil
		},
	}

	cmd.Flags().StringP("host", "H", "", "Database host")
	cmd.Flags().IntP("port", "P", 0, "Database port")
	cmd.Flags().StringP("dbname", "N", "", "Database name")
	cmd.Flags().StringP("username", "U", "", "Username")
	cmd.Flags().StringP("password", "p", "", "Password")

	return cmd
}
