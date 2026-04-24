package commands

import (
	"dbtool/DBdrivers"
	"fmt"

	"github.com/spf13/cobra"
)

func BackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "backup",
		Short:        "command for create backup db",
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
			conf.FilePath = GlobalViper.GetString("path")
			if conf.Host == "" {
				return fmt.Errorf("host is empty. Set it in config.json or use --host flag")
			}
			if conf.User == "" {
				return fmt.Errorf("username is empty. Set it in config.json or use --username flag")
			}
			if conf.Port == 0 {
				conf.Port = 5432
			}
			if RuntimeErr := DBdrivers.BackupPostgres(conf); RuntimeErr != nil {
				return RuntimeErr
			}
			return nil
		},
	}
	cmd.Flags().StringP("host", "H", "", "Host")
	cmd.Flags().IntP("port", "P", 0, "Port")
	cmd.Flags().StringP("database", "", "", "Type DB")
	cmd.Flags().StringP("dbname", "N", "", "Name db to connect")
	cmd.Flags().StringP("username", "U", "", "Username")
	cmd.Flags().StringP("password", "", "", "Password")
	cmd.Flags().StringP("path", "", ".", "The path where the file will be saved")
	return cmd
}
