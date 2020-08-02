package api

import (
	"fmt"
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/luqmansen/web-analytics/configs"
	"github.com/spf13/viper"
	"net/http"
)

func HostValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conf := configs.GetConfig()
		if r.Method == "POST" {
			if found := func() bool {
				for _, i := range conf.AllowedRequest {
					fmt.Println(i, r.Host)
					if i == r.Host {
						return true
					}
				}
				return false
			}(); !found {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Error: " + analytics.ErrorInvalidHost.Error()))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conf := configs.GetConfig()

		if viper.GetString("DEPLOY") != "PROD" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Add("Access-Control-Allow-Origin", conf.AllowedHost)
		}

		w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
