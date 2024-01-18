package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string `mapstructure:"CLCK_ENV"`
	Port       string `mapstructure:"CLCK_PORT"`
	DBServer   string `mapstructure:"CLCK_DB_SERVER"`
	DBName     string `mapstructure:"CLCK_DB_NAME"`
	DBUser     string `mapstructure:"CLCK_DB_USER"`
	DBPassword string `mapstructure:"CLCK_DB_PASS"`
}

func MustLoad(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	config.Port = "9090"
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	err = viper.Unmarshal(&config)
	return
}
