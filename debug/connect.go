package debug

import (
	"crypto/tls"
	"net"
)

func ConnectTCP(remoteAddr string) (err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	_ = conn.Close()

	return
}

func ConnectTLS(remoteAddr string, insecureSkipVerify bool) (err error) {
	conf := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify, // nolint: gosec
	}

	conn, err := tls.Dial("tcp", remoteAddr, conf)
	if err != nil {
		return
	}

	_ = conn.Close()

	return
}
