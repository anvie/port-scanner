package predictors

import (
//	"net"
)

type Predictor interface {
	Predict(host string) string
}

