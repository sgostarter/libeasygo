package grpctoolset

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/sgostarter/libeasygo/commerr"
)

func GenServerTLSConfig(cfg *GRPCTlsConfig) (tlsConfig *tls.Config, err error) {
	if cfg == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	caPool := x509.NewCertPool()

	for _, ca := range cfg.RootCAs {
		caPool.AppendCertsFromPEM(ca)
	}

	cert, err := tls.X509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return
	}

	// nolint: gosec
	tlsConfig = &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
	}

	return
}

func GenClientTLSConfig(cfg *GRPCTlsConfig) (tlsConfig *tls.Config, err error) {
	if cfg == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	caPool := x509.NewCertPool()

	for _, ca := range cfg.RootCAs {
		caPool.AppendCertsFromPEM(ca)
	}

	cert, err := tls.X509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return
	}

	// nolint: gosec
	tlsConfig = &tls.Config{
		ServerName:   cfg.ServerName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
	}

	return
}
