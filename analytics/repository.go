package analytics

type AnalyticRepository interface {
	GetAll() ([]*Analytic, error)
	Store(analytic *Analytic) error
}
