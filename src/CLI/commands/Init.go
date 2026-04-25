package commands

import (
	"dbtool/DBinterface"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GlobalViper *viper.Viper

func InitConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	v.SetEnvPrefix("DB")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		fmt.Println("Config file not found")
	}

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	GlobalViper = v
	return nil
}

func InitCmd(cmd *cobra.Command, conf *DBinterface.Config) error {

	if err := InitConfig(cmd); err != nil {
		return fmt.Errorf("init config error: %w", err)
	}

	conf.Host = GlobalViper.GetString("host")
	conf.Port = GlobalViper.GetInt("port")
	conf.User = GlobalViper.GetString("username")
	conf.Password = GlobalViper.GetString("password")
	conf.DBName = GlobalViper.GetString("dbname")
	conf.DBtype = GlobalViper.GetString("database")

	if conf.DBtype == "" {
		return fmt.Errorf("database type is empty. Set it in config.json or use --host flag")
	}
	if conf.Host == "" {
		return fmt.Errorf("host is empty. Set it in config.json or use --host flag")
	}
	if conf.User == "" {
		return fmt.Errorf("username is empty. Set it in config.json or use --username flag")
	}
	if conf.Port == 0 {
		return fmt.Errorf("port is empty. Set it in config.json or use --username flag")

	}

	return nil
}
