package netutils

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"time"
)

type TLSConfigModifier func(cfg *tls.Config)

func TLSConfigModifier4InsecureSkipVerify(cfg *tls.Config) {
	if cfg == nil {
		return
	}

	cfg.InsecureSkipVerify = true
}

func DialTCPWithTimeout(ctx context.Context, useSSL bool, address string, timeout time.Duration,
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

	/*
		d := tls.Dialer{
			Config: tlsConfig,
		}

		return d.DialContext(ctx, "tcp", address)
	*/

	c, err := dial(ctx, new(net.Dialer), "tcp", address, tlsConfig)
	if err != nil {
		// Don't return c (a typed nil) in an interface.
		return nil, err
	}

	return c, nil
}

//
// go\src\crypto\tls\tls.go
// to fix
//

var emptyConfig tls.Config

func defaultConfig() *tls.Config {
	return &emptyConfig
}

// nolint
func dial(ctx context.Context, netDialer *net.Dialer, network, addr string, config *tls.Config) (*tls.Conn, error) {
	if netDialer.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, netDialer.Timeout)
		defer cancel()
	}

	if !netDialer.Deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, netDialer.Deadline)
		defer cancel()
	}

	rawConn, err := netDialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	colonPos := strings.LastIndex(addr, ":")
	if colonPos == -1 {
		colonPos = len(addr)
	}
	hostname := addr[:colonPos]

	if config == nil {
		config = defaultConfig()
	}

	// If no ServerName is set, infer the ServerName
	// from the hostname we're connecting to.
	if !config.InsecureSkipVerify {
		if config.ServerName == "" {
			// Make a copy to avoid polluting argument or default.
			c := config.Clone()
			c.ServerName = hostname
			config = c
		}
	} else {
		if config.ServerName != "" {
			// Make a copy to avoid polluting argument or default.
			c := config.Clone()
			c.ServerName = ""
			config = c
		}
	}

	conn := tls.Client(rawConn, config)

	err = conn.HandshakeContext(ctx)
	if err != nil {
		rawConn.Close()

		return nil, err
	}

	return conn, nil
}
