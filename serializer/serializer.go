package serializer

import "github.com/luqmansen/web-analytics/analytics"

type AnalyticSerializer interface {
	Decode(input []byte) (*analytics.Analytic, error)
	Encode(input *analytics.Analytic) ([]byte, error)
}
