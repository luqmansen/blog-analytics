package api

import (
	"github.com/spf13/viper"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if viper.GetString("DEPLOY") != "PROD" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Add("Access-Control-Allow-Origin", viper.GetString("AllowedHost"))
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
