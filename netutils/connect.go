package netutils

import (
	"context"
	"net"
	"time"

	"github.com/sgostarter/i/commerr"
)

func dialTCPEx(useSSL, ignoreCert bool, address string, timeout time.Duration) (net.Conn, error) {
	var modifiers []TLSConfigModifier

	if useSSL && ignoreCert {
		modifiers = append(modifiers, TLSConfigModifier4InsecureSkipVerify)
	}

	return DialTCPWithTimeout(context.TODO(), useSSL, address, timeout, modifiers...)
}

func TestTCPConnect(remoteAddr string, useTLS bool) (err error) {
	return TestTCPConnectEx(remoteAddr, useTLS, time.Second*10, false)
}

func TestTCPConnectEx(remoteAddr string, useTLS bool, timeout time.Duration, strictMode bool) (err error) {
	err = TestTCPConnectWithTimeout(remoteAddr, useTLS, timeout)
	if err != nil {
		return
	}

	if !strictMode || useTLS {
		return
	}

	if TestTCPConnectWithTimeout(remoteAddr, true, timeout) == nil {
		err = commerr.ErrReject
	}

	return
}

func TestTCPConnectWithTimeout(remoteAddr string, useTLS bool, timeout time.Duration) (err error) {
	return TestTCPConnectWithTimeoutEx(remoteAddr, useTLS, true, timeout)
}

func TestTCPConnectWithTimeoutEx(remoteAddr string, useTLS, ignoreCert bool, timeout time.Duration) (err error) {
	conn, err := dialTCPEx(useTLS, ignoreCert, remoteAddr, timeout)
	if err != nil {
		return
	}

	_ = conn.Close()

	return
}
