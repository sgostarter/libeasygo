package netutils

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	utListenAddress = "127.0.0.1:23455"
)

func utReadFileContent(t *testing.T, file string) []byte {
	d, err := os.ReadFile(file)
	assert.Nil(t, err)

	return d
}

func TestListener(t *testing.T) {
	t.SkipNow()

	certsRoot := "../certs/certs_svc"

	var tlsCertificate tls.Certificate

	tlsCertificate, err := tls.X509KeyPair(utReadFileContent(t, filepath.Join(certsRoot, "server_cert.pem")),
		utReadFileContent(t, filepath.Join(certsRoot, "server_key.pem")))
	assert.Nil(t, err)

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(utReadFileContent(t, filepath.Join(certsRoot, "client_ca_cert.pem"))); !ok {
		t.Log("No certs appended, using system certs only")
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{tlsCertificate},
		ClientCAs:    rootCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}
	tlsCfg.Rand = rand.Reader

	listener, err := NewTLSAndTCPTransport(utListenAddress, tlsCfg)
	assert.Nil(t, err)

	c, err := listener.Accept()
	assert.Nil(t, err)

	r := bufio.NewReader(c)
	s, err := r.ReadString('\n')
	assert.Nil(t, err)
	t.Log(s)

	w := bufio.NewWriter(c)
	_, err = w.WriteString("hello, client\n")
	assert.Nil(t, err)

	_ = w.Flush()

	_ = c.Close()
}

func TestTCPClient(t *testing.T) {
	t.SkipNow()

	c, err := net.Dial("tcp", utListenAddress)
	assert.Nil(t, err)

	defer c.Close()

	w := bufio.NewWriter(c)
	_, err = w.WriteString("hello, server\n")
	assert.Nil(t, err)

	_ = w.Flush()

	r := bufio.NewReader(c)
	s, err := r.ReadString('\n')
	assert.Nil(t, err)
	t.Log(s)
}

func TestTLSClient(t *testing.T) {
	t.SkipNow()

	certsRoot := "../certs/certs_cli"

	var tlsCertificate tls.Certificate

	tlsCertificate, err := tls.X509KeyPair(utReadFileContent(t, filepath.Join(certsRoot, "client_cert.pem")),
		utReadFileContent(t, filepath.Join(certsRoot, "client_key.pem")))
	assert.Nil(t, err)

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(utReadFileContent(t, filepath.Join(certsRoot, "server_ca_cert.pem"))); !ok {
		t.Log("No certs appended, using system certs only")
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{
		tlsCertificate},
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	}
	tlsCfg.Rand = rand.Reader

	c, err := tls.Dial("tcp", utListenAddress, tlsCfg)
	assert.Nil(t, err)

	version := c.ConnectionState().Version
	t.Logf("Selected TLS version: %d", version)
	t.Logf("Selected TLS version: 0x%04x", version)

	w := bufio.NewWriter(c)
	_, err = w.WriteString("hello, server\n")
	assert.Nil(t, err)

	_ = w.Flush()

	r := bufio.NewReader(c)
	s, err := r.ReadString('\n')
	assert.Nil(t, err)
	t.Log(s)
}
