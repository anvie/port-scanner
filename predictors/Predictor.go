package predictors

import (
//	"net"
	"strings"
)

type Predictor interface {
	DetailPredictor
	Predict(host string) string
	PredictResponse(resp string, dp DetailPredictor) string
}

type DetailPredictor interface {
	PredictResponseDetail(resp string) string
}

type BaseHttpPredictor struct {
//	DetailPredictor
}

func (pa *BaseHttpPredictor) PredictResponse(resp string, dp DetailPredictor) string {
	if strings.Contains(resp, "HTTP/") {
		rv := ""
		detail := dp.PredictResponseDetail(resp)
		if len(detail) > 0 {
			rv = "web server"
		}
		return strings.TrimSpace(rv + " " + detail)
	}
	return ""
}

func (pa *BaseHttpPredictor) PredictResponseDetail(resp string) string {
	return "aoeu"
}
