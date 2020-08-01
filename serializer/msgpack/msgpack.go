package msgpack

import (
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type Analytic struct{}

func (a *Analytic) Decode(input []byte) (*analytics.Analytic, error) {
	analytic := &analytics.Analytic{}
	if err := msgpack.Unmarshal(input, analytic); err != nil {
		return nil, errors.Wrap(err, "serializer.msgpack.Decode")
	}
	return analytic, nil
}

func (a *Analytic) Encode(input *analytics.Analytic) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.msgpack.Encode")
	}
	return rawMsg, nil
}
