/***
 * Simple port scanner
 *
 *
 */

package portscanner

import (
	"net"
	//	"os"
	"fmt"
	//	"io/ioutil"
	//	"strings"
	"time"
	"sync"

	"github.com/anvie/port-scanner/predictors"
	"github.com/anvie/port-scanner/predictors/webserver"
)

type PortScanner struct {
	host       string
	predictors []predictors.Predictor
	timeout    time.Duration
	threads    int
	usePredictor bool
}

func NewPortScanner(host string, timeout time.Duration, threads int) *PortScanner {
	return &PortScanner{host, []predictors.Predictor{
		&webserver.ApachePredictor{},
		&webserver.NginxPredictor{},
	}, timeout, threads, true,
	}
}
func (h *PortScanner) TogglePredictor(usePredictor bool) {
	h.usePredictor = usePredictor
}
func (h *PortScanner) SetThreads(threads int) {
	h.threads = threads
}
func (h *PortScanner) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
}
func (h *PortScanner) RegisterPredictor(predictor predictors.Predictor) {
	for _, p := range h.predictors {
		if p == predictor {
			return
		}
	}
	h.predictors = append(h.predictors, predictor)
}

func (h PortScanner) IsOpen(port int) bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", h.hostPort(port))
	if err != nil {
		return false
	}
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), h.timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}

func (h PortScanner) GetOpenedPort(portStart int, portEnds int) []int {
	rv := []int{}
	l := sync.Mutex{}
	sem := make(chan bool, h.threads)
	for port := portStart; port <= portEnds; port++ {
		sem <- true
		go func(port int) {
			if h.IsOpen(port) {
				l.Lock()
				rv = append(rv, port)
				l.Unlock()
			}
			<- sem
		}(port)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return rv
}

func (h PortScanner) hostPort(port int) string {
	return fmt.Sprintf("%s:%d", h.host, port)
}

const UNKNOWN = "<unknown>"

func (h PortScanner) openConn(host string) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout("tcp", tcpAddr.String(), h.timeout)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (h PortScanner) DescribePort(port int) string {
	if !h.usePredictor {
		return h.predictPort(port)
	}

	switch {
	default:
		return UNKNOWN
	case h.IsHttp(port):
		rv := h.PredictUsingPredictor(h.hostPort(port))
		return rv
	case port > 0:
		assumed := h.predictPort(port)
		rv := assumed
		if assumed == UNKNOWN {
			rv = h.PredictUsingPredictor(h.hostPort(port))
		}

		switch assumed {
		case "MySQL":
			// get the version
			conn, err := h.openConn(h.hostPort(port))
			if err == nil {
				defer conn.Close()

				duration, _ := time.ParseDuration("3s")

				conn.SetDeadline(time.Now().Add(duration))

				conn.Read(make([]byte, 5))
				result := make([]byte, 20)

				_, err := conn.Read(result)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return ""
				}

				resp := string(result)
				rv = assumed + " version: " + resp
			}
		}

		return rv
	}
}

func (h PortScanner) IsHttp(port int) bool {
	return port == 80 || port == 8080
}

func (h PortScanner) PredictUsingPredictor(host string) string {
	for _, p := range h.predictors {
		conn, err := h.openConn(host)
		if err != nil {
			break
		}
		defer conn.Close()
		rv := p.Predict(host)
		if len(rv) > 0 {
			return rv
		}
	}
	return UNKNOWN
}

var KNOWN_PORTS = map[int]string{
	27017: "mongodb [ http://www.mongodb.org/ ]",
	28017: "mongodb web admin [ http://www.mongodb.org/ ]",
	21:    "ftp",
	22:    "SSH",
	23:    "telnet",
	25:    "SMTP",
	66:    "Oracle SQL*NET?",
	69:    "tftp",
	80:    "http",
	88:    "kerberos",
	109:   "pop2",
	110:   "pop3",
	123:   "ntp",
	137:   "netbios",
	139:   "netbios",
	443:   "https",
	445:   "Samba",
	631:   "cups",
	5800:  "VNC remote desktop",
	194:   "IRC",
	118:   "SQL service?",
	150:   "SQL-net?",
	1433:  "Microsoft SQL server",
	1434:  "Microsoft SQL monitor",
	3306:  "MySQL",
	3396:  "Novell NDPS Printer Agent",
	3535:  "SMTP (alternate)",
	554:   "RTSP",
	9160:  "Cassandra [ http://cassandra.apache.org/ ]",
}

func (h PortScanner) predictPort(port int) string {
	if rv, ok := KNOWN_PORTS[port]; ok {
		return rv
	}
	return UNKNOWN
}
