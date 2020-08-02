package analytics

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrorInvalidURL  = errors.New("Invalid URI.")
	ErrorInvalidHost = errors.New("Host not allowed.")
	ErrorDuplicate   = errors.New("Duplicate data.")
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

func (a analyticServices) Store(analytic *Analytic) error {
	if err := a.Validate(analytic); err != nil {
		return err
	}
	analytic.CreatedAt = time.Now().UTC()
	return a.analyticRepo.Store(analytic)
}

func (a analyticServices) Validate(analytic *Analytic) error {

	if err := validate.Validate(*analytic); err != nil {
		logrus.Debugln(*analytic)
		return errors.Wrap(err, "service.Analytic.Validate")
	}
	return nil
}
