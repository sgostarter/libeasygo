package netutils

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

func DialTCPWithTimeout(ctx context.Context, useSSL bool, address string, timeout time.Duration) (net.Conn, error) {
	return DialTCPWithTimeoutEx(ctx, useSSL, false, address, timeout)
}

func DialTCPWithTimeoutEx(ctx context.Context, useSSL, insecureSkipVerify bool, address string, timeout time.Duration) (net.Conn, error) {
	if !useSSL {
		return net.DialTimeout("tcp", address, timeout)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	d := tls.Dialer{
		// nolint: gosec
		Config: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
	}

	return d.DialContext(ctx, "tcp", address)
}
