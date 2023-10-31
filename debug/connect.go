package debug

import (
	"context"
	"net"
	"time"

	"github.com/sgostarter/libeasygo/netutils"
)

func dialTCPEx(useSSL bool, address string, timeout time.Duration) (net.Conn, error) {
	var modifiers []netutils.TLSConfigModifier

	if useSSL {
		modifiers = append(modifiers, netutils.TLSConfigModifier4InsecureSkipVerify)
	}

	return netutils.DialTCPWithTimeout(context.TODO(), useSSL, address, timeout, modifiers...)
}

func TestTCPConnect(remoteAddr string, useTLS bool) (err error) {
	return TestTCPConnectWithTimeout(remoteAddr, useTLS, time.Second*10)
}

func TestTCPConnectWithTimeout(remoteAddr string, useTLS bool, timeout time.Duration) (err error) {
	conn, err := dialTCPEx(useTLS, remoteAddr, timeout)
	if err != nil {
		return
	}

	_ = conn.Close()

	return
}
