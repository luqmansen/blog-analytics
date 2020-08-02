package api

import (
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"net/url"
)

func HostValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowed := viper.GetStringSlice("AllowedRequest")

		if r.Method == "POST" {
			if found := func() bool {
				for _, i := range allowed {
					// This is absolutely disgusting, idk what better. I want the configuration
					// should be able to work with arbitrary port. The net.SplitHostPort can't
					// detect if Hostname doesn't have port, it'll instead detect u scheme as
					//the host. eg http://localhost -> http returned as host

					u, _ := url.Parse(i)
					h1, _, err := net.SplitHostPort(u.Host)
					h2, _, _ := net.SplitHostPort(r.Host)

					if err != nil{
						h1 = u.Host
					}

					if h1 == h2 {
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
