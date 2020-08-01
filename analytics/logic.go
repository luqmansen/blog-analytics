package analytics

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/dealancer/validate.v2"
	url2 "net/url"
	"time"
)

var (
	ErrorInvalidURL  = errors.New("Invalid URI")
	ErrorInvalidHost = errors.New("Invalid Host")
)

type analyticServices struct {
	analyticRepo AnalyticRepository
}

func NewAnalyticService(repository AnalyticRepository) AnalyticServices {
	return &analyticServices{
		repository,
	}
}

func (a analyticServices) GetAll() ([]*Analytic, error) {
	return a.analyticRepo.GetAll()
}

func (a analyticServices) Validate(analytic *Analytic) error {
	if err := validate.Validate(analytic); err != nil {
		return errors.Wrap(ErrorInvalidURL, "service.Analytic.Validate")
	}
	url, err := url2.Parse(analytic.URL)
	if err != nil {
		return err
	}
	t := viper.GetStringSlice("validClient")
	for _, v := range t {
		if url.Host != v {
			return errors.Wrap(ErrorInvalidHost, "service.Analytic.Validate")
		}
	}
	return nil
}

func (a analyticServices) Store(analytic *Analytic) error {
	if err := a.Validate(analytic); err != nil {
		return err
	}
	analytic.CreatedAt = time.Now().UTC().Unix()
	return a.analyticRepo.Store(analytic)
}
