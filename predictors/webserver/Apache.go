package webserver

import (
	"net"
	"strings"
	"io/ioutil"
	"time"
)

type ApachePredictor struct {

}


func (p ApachePredictor) Predict(host string) string {
	duration, _ := time.ParseDuration("10s")

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

	resp := strings.ToLower(string(result))
	if strings.Contains(resp, "HTTP/") {
		rv := "web service"
		if strings.Contains(resp, "apache") {
			rv = rv + " Apache"
		}
		return rv
	}
	return ""
}

