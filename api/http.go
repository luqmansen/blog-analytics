package api

import (
	"encoding/json"
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/luqmansen/web-analytics/serializer"
	js "github.com/luqmansen/web-analytics/serializer/json"
	ms "github.com/luqmansen/web-analytics/serializer/msgpack"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

type AnalyticsHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	analyticsService analytics.AnalyticServices
}

func NewHandler(services analytics.AnalyticServices) AnalyticsHandler {
	return &handler{analyticsService: services}
}

func setResponse(w http.ResponseWriter, contentType string, body []byte, status int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) serializer.AnalyticSerializer {
	if contentType == "application/x-msgpack" {
		return &ms.Analytic{}
	}
	return &js.Analytic{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	all, err := h.analyticsService.GetAll()
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
	}
	//dict := map[string]interface{}{"data": all}
	b, err := json.Marshal(all)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
	}

	w.Write(b)
	w.WriteHeader(200)
}

func (h handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	analytic, err := h.serializer(contentType).Decode(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.analyticsService.Store(analytic)
	if err != nil {
		if errors.Cause(err) == analytics.ErrorInvalidURL {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
