package commands

import (
	"dbtool/Cloud"
	"dbtool/DBinterface"
	"dbtool/DBinterface/DBdrivers"
	"fmt"

	"github.com/spf13/cobra"
)

func BackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "backup",
		Short:        "command for create backup db",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			conf := DBinterface.Config{}
			if err := InitCmd(cmd, &conf); err != nil {
				return err
			}
			conf.FilePath = GlobalViper.GetString("path")

			CloudConf := Cloud.CloudConfig{}
			if CloudConf.IsUse = GlobalViper.GetBool("usecloud"); CloudConf.IsUse {
				if err := Cloud.InitCloud(&CloudConf, GlobalViper); err != nil {
					return err
				}
			}

			switch conf.DBtype {
			case "postgres":
				{
					if RuntimeErr := DBdrivers.BackupPostgres(conf); RuntimeErr != nil {
						return RuntimeErr
					}
					if CloudConf.IsUse {
						if RuntimeErr := Cloud.BackupCloud(CloudConf, conf); RuntimeErr != nil {
							return RuntimeErr
						}
					}
				}
			case "mysql":
				{
					if RuntimeErr := DBdrivers.BackupMySQL(conf); RuntimeErr != nil {
						return RuntimeErr
					}
					if CloudConf.IsUse {
						if RuntimeErr := Cloud.BackupCloud(CloudConf, conf); RuntimeErr != nil {
							return RuntimeErr
						}
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
	cmd.Flags().StringP("database", "", "", "Type DB")
	cmd.Flags().StringP("dbname", "N", "", "Name db to connect")
	cmd.Flags().StringP("username", "U", "", "Username")
	cmd.Flags().StringP("password", "", "", "Password")
	cmd.Flags().StringP("path", "", ".", "The path where the file will be saved")
	cmd.Flags().BoolP("usecloud", "", false, "Use cloud or not?")
	cmd.Flags().StringP("accesskey", "A", "", "Access key for S3")
	cmd.Flags().StringP("secretkey", "S", "", "Secret key for S3")
	cmd.Flags().StringP("endpoint", "E", "", "Endpoint")
	cmd.Flags().StringP("bucketname", "B", "", "Name of bucket where file is")
	return cmd
}
