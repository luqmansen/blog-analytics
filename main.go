package main

import (
	"fmt"
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/luqmansen/web-analytics/api"
	"github.com/luqmansen/web-analytics/configs"
	_ "github.com/luqmansen/web-analytics/configs"
	"github.com/luqmansen/web-analytics/repository/mongodb"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	config := configs.GetConfig()

	repo, err := mongodb.NewMongoRepository(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	service := analytics.NewAnalyticService(repo)
	handler := api.NewHandler(service)

	r := api.Routes(handler)

	errs := make(chan error, 2)
	go func() {
		logrus.Infof("Listening on http://%s:%s", config.Server.Host, config.Server.Port)
		errs <- http.ListenAndServe(":"+config.Server.Port, r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Infof("Terminated %s", <-errs)
}
