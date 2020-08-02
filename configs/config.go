package configs

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Configuration struct {
	Server         Server
	Database       Database
}

type Database struct {
	URI      string
	Database string
	Timeout  int
}

type Server struct {
	Host string
	Port string
}

var config Configuration

func init() {
	viper.SetDefault("DEPLOY", "PROD")

	if os.Getenv("DEPLOY") != "PROD" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8080")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logrus.Fatal(fmt.Errorf("Fatal error failed to decode to struct: %s \n", err))
	}

	if len(viper.GetStringSlice("AllowedRequest")) == 0 {
		logrus.Errorln(errors.New("valid request host not set"))
	}
	if viper.GetString("DEPLOY") == "PROD" {
		config.Server.Port = viper.GetString("PORT")
		config.Server.Host = viper.GetString("HOST")
	}

}

func GetConfig() *Configuration {
	return &config
}
