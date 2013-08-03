package webserver

import (
	"net"
	"io/ioutil"
	"strings"
	"time"
	"github.com/anvie/port-scanner/predictors"
//	"fmt"
)

type NginxPredictor struct {
	predictors.BaseHttpPredictor
}


func (p *NginxPredictor) Predict(host string) string {
	duration, _ := time.ParseDuration("3s")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if (err != nil) {
		return ""
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return ""
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(duration))

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	if err != nil {
		return ""
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return ""
	}

	resp := string(result)
	return p.PredictResponse(resp, p)
}

func (p *NginxPredictor) PredictResponseDetail(resp string) string {
	if strings.Contains(resp, "nginx/") {
		return "Nginx [ http://wiki.nginx.org/Main ]"
	}
	return ""
}


