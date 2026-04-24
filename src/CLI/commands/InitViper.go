package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GlobalViper *viper.Viper

func InitConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("./CLI/")

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
