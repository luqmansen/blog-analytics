package main

import (
	"fmt"

	"github.com/luqmansen/web-analytics/analytics"
	"github.com/luqmansen/web-analytics/api"
	"github.com/luqmansen/web-analytics/configs"
	"github.com/luqmansen/web-analytics/repository/mongodb"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var config configs.Configuration

func init() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error failed to decode to struct: %s \n", err))
	}

}

func main() {
	repo, err := mongodb.NewMongoRepository(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	service := analytics.NewAnalyticService(repo)
	handler := api.NewHandler(service)

	r := api.Routes(handler)

	errs := make(chan error, 2)
	go func() {
		fmt.Printf("Listening on port %s\n", config.Server.Port)
		errs <- http.ListenAndServe(config.Server.Port, r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}
