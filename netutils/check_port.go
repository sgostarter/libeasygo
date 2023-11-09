package netutils

import (
	"fmt"
	"net"
)

func GetAvailableTCPPort() (freePort int, err error) {
	return GetAvailableTCPPortEx(0)
}

func GetAvailableTCPPortEx(maxRetry int) (freePort int, err error) {
	if maxRetry <= 0 {
		maxRetry = 0
	}

	for ; maxRetry >= 0; maxRetry-- {
		freePort, err = getAvailableTCPPortEx()
		if err == nil {
			break
		}
	}

	return
}

func getAvailableTCPPortEx() (freePort int, err error) {
	address, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return
	}

	if tcpAddr, ok := listener.Addr().(*net.TCPAddr); ok {
		freePort = tcpAddr.Port
	}

	_ = listener.Close()

	return
}

func CheckTCPPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}

	_ = listener.Close()

	return true
}
