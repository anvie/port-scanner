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
	"io/ioutil"
	"strings"
)


type PortScanner struct {
	host string
}

func NewPortScanner(host string) *PortScanner {
	return &PortScanner{host}
}

func (h PortScanner) IsOpen(port int) bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", h.hostPort(port))
	if (err != nil) {
		return false
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if (err != nil) {
		return false
	}

	defer conn.Close()

	return true
}


func (h PortScanner) GetOpenedPort(portStart int, portEnds int) []int {
	rv := []int{}
	for port := portStart; port <= portEnds; port++ {
		if h.IsOpen(port) {
			rv = append(rv, port)
		}
	}
	return rv
}

func (h PortScanner) hostPort(port int) string {
	return fmt.Sprintf("%s:%d", h.host, port)
}

const UNKNOWN = "<unknown>"

func (h PortScanner) DescribePort(port int) string {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", h.hostPort(port))
	if (err != nil) {
		return UNKNOWN
	}
	switch {
	default:
		return UNKNOWN
	case h.IsHttp(port):
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return UNKNOWN
		}
		defer conn.Close()

		_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		if err != nil {
			return UNKNOWN
		}

		result, err := ioutil.ReadAll(conn)
		if err != nil {
			return UNKNOWN
		}

		resp := string(result)
		rv := h.predictResponse(resp)

		return rv
	case port > 0:
		return h.predictPort(port)
	} 
}


func (h PortScanner) IsHttp(port int) bool {
	return port == 80 || port == 8080
}

func (h PortScanner) predictResponse(resp string) string {
	if strings.Contains(resp, "HTTP/") {
		rv := "web service"
		if strings.Contains(resp, "nginx/") {
			rv = rv + " nginx"
		}
		return rv
	}
	return UNKNOWN
}

var KNOWN_PORTS = map[int]string {
	27017: "mongodb [ http://www.mongodb.org/ ]",
	28017: "mongodb web admin [ http://www.mongodb.org/ ]",
	21: "ftp",
	22: "SSH",
	23: "telnet",
	25: "SMTP",
	66: "Oracle SQL*NET?",
	69: "tftp",
	88: "kerberos",
	109: "pop2",
	110: "pop3",
	123: "ntp",
	137: "netbios",
	139: "netbios",
	445: "Samba",
	631: "cups",
	5800: "VNC remote desktop",
	194: "IRC",
	118: "SQL service?",
	150: "SQL-net?",
	1433: "Microsoft SQL server",
	1434: "Microsoft SQL monitor",
	3306: "MySQL",
	3396: "Novell NDPS Printer Agent",
	3535: "SMTP (alternate)",
	554: "RTSP",
}

func (h PortScanner) predictPort(port int) string {
	if rv, ok := KNOWN_PORTS[port]; ok {
		return rv
	}
	return UNKNOWN
}

