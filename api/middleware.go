package api

import (
	"fmt"
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/luqmansen/web-analytics/configs"
	"net/http"
)

func HostValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conf := configs.GetConfig()
		if r.Method == "POST" {
			if found := func() bool {
				for _, i := range conf.AllowedHosts {
					fmt.Printf("%s %s\n", i, r.Host)
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
