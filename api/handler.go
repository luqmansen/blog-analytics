package api

import (
	"encoding/json"
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/luqmansen/web-analytics/serializer"
	js "github.com/luqmansen/web-analytics/serializer/json"
	ms "github.com/luqmansen/web-analytics/serializer/msgpack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
		logrus.Error(err)
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
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if len(all) == 0 {
		w.Write(nil)
		w.WriteHeader(http.StatusNoContent)
	}

	b, err := json.Marshal(all)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
	w.WriteHeader(200)
}

func (h handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Empty body not allowed", http.StatusBadRequest)
		return
	}

	analytic, err := h.serializer(contentType).Decode(body)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	analytic.IP = r.RemoteAddr
	err = h.analyticsService.Store(analytic)
	if err != nil {
		cause := errors.Cause(err)
		if (cause == analytics.ErrorDuplicate) || (cause == analytics.ErrorInvalidURL) || (cause == analytics.ErrorInvalidHost) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorln(err)
		return
	}
	w.Write([]byte("Store analytics success"))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Health Check")
	w.Write([]byte("OK"))
	w.WriteHeader(200)
}
