package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GlobalViper *viper.Viper

func InitConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("JSON")
	v.AddConfigPath(".")

	v.SetEnvPrefix("DB")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	GlobalViper = v
	return nil
}
