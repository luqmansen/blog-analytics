package analytics

import (
	"cloud.google.com/go/civil"
	"github.com/pkg/errors"
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
	// this is seems redundant, but currently stdlib
	// not support time.Date field
	analytic.CreatedAt = civil.DateTime{
		Date: civil.DateOf(time.Now().UTC()),
		Time: civil.TimeOf(time.Now().UTC()),
	}
	return a.analyticRepo.Store(analytic)
}

func (a analyticServices) Validate(analytic *Analytic) error {

	if err := validate.Validate(*analytic); err != nil {
		return errors.Wrap(err, "service.Analytic.Validate")
	}
	return nil
}
