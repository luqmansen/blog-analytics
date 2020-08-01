package json

import (
	"encoding/json"
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/pkg/errors"
)

type Analytic struct{}

func (a *Analytic) Decode(input []byte) (*analytics.Analytic, error) {
	analytic := &analytics.Analytic{}
	if err := json.Unmarshal(input, analytic); err != nil {
		return nil, errors.Wrap(err, "serializer.json.Decode")
	}
	return analytic, nil
}

func (a *Analytic) Encode(input *analytics.Analytic) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.json.Encode")
	}
	return rawMsg, nil
}
