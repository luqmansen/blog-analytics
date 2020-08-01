package analytics

type AnalyticServices interface {
	Validate(analytic *Analytic) error
	GetAll() ([]*Analytic, error)
	Store(analytic *Analytic) error
}
