package netutils

import (
	"bufio"
	"crypto/tls"
	"net"
)

func NewTLSAndTCPTransport(listenAddress string, tlsConfig *tls.Config) (listener net.Listener, err error) {
	tcpListener, err := net.Listen("tcp4", listenAddress)
	if err != nil {
		return
	}

	listener = &customListener{
		Listener:  tcpListener,
		tlsConfig: tlsConfig,
	}

	return
}

type customListener struct {
	net.Listener
	tlsConfig *tls.Config
}

func (li *customListener) Accept() (net.Conn, error) {
	c, err := li.Listener.Accept()
	if err != nil {
		return nil, err
	}

	buffConn := NewBufferedConn(c)

	b, err := buffConn.Peek(6)
	if err != nil {
		return nil, err
	}

	if IsTLS(b) {
		return tls.Server(buffConn, li.tlsConfig), nil
	}

	return buffConn, nil
}

// nolint
/*
record type: 1 byte (0x16 for "records contains some handshake message data")
protocol version: 2 bytes (0x03 0x00 for SSL 3.0, 0x03 0x01 for TLS 1.0, and so on)
record length: 2 bytes (big endian)
then the record data...
For the first record (from client to server), the client will first send a ClientHello message which is a type of handshake message, hence encapsulated in a record as shown above (the first byte of the record will be 0x16). Theoretically, the client may send the ClientHello split into several records, and it may begin with one or several empty records, but this is not very probable. The ClientHello message itself begins with its own four-byte header, with one byte for the message type (0x01 for ClientHello), then the message length over three bytes (there again, big-endian).

Once the client has sent its ClientHello, then it expects a response from the server, so the ClientHello will be alone in its record.

So you could expect a payload which begins with the following 9 bytes:

0x16 0x03 X Y Z 0x01 A B C
with:

X will be 0, 1, 2, 3... or more, depending on the protocol version used by the client for this first message. Currently, defined SSL/TLS versions are SSL 3.0, TLS 1.0, TLS 1.1 and TLS 1.2. Other versions may be defined in the future. They will probably use the 3.X numbering scheme, so you can expect the second header byte to remain a 0x03, but you should not arbitrarily limit the third byte.

Y Z is the encoding of the record length; A B C is the encoding of the ClientHello message length. Since the ClientHello message begins with a 4-byte header (not including in its length) and is supposed to be alone in its record, you should have: A = 0 and 256*X+Y = 256*B+C+4.

If you see 9 such bytes, which verify these conditions, then chances are that this is a ClientHello from a SSL client.
*/

func IsTLS(b []byte) bool {
	return b[0] == 0x16 && b[5] == 0x01
}

type BufferedConn struct {
	r *bufio.Reader
	net.Conn
}

func NewBufferedConn(conn net.Conn) *BufferedConn {
	return &BufferedConn{
		r:    bufio.NewReader(conn),
		Conn: conn,
	}
}

func (b *BufferedConn) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b *BufferedConn) Read(p []byte) (int, error) {
	return b.r.Read(p)
}
