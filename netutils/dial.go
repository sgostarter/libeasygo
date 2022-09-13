package netutils

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

type TLSConfigModifier func(cfg *tls.Config)

func TLSConfigModifier4InsecureSkipVerify(cfg *tls.Config) {
	if cfg == nil {
		return
	}

	cfg.InsecureSkipVerify = true
}

func DialTCPWithTimeout(ctx context.Context, useSSL bool, address string, timeout time.Duration) (net.Conn, error) {
	return DialTCPWithTimeoutEx(ctx, useSSL, address, timeout)
}

func DialTCPWithTimeoutEx(ctx context.Context, useSSL bool, address string, timeout time.Duration,
	tlsConfigModifier ...TLSConfigModifier) (net.Conn, error) {
	if !useSSL {
		return net.DialTimeout("tcp", address, timeout)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// nolint: gosec
	tlsConfig := &tls.Config{}
	for _, modifier := range tlsConfigModifier {
		modifier(tlsConfig)
	}

	d := tls.Dialer{
		Config: tlsConfig,
	}

	return d.DialContext(ctx, "tcp", address)
}
