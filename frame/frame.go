package frame

type Frame interface {
	WriteFrame(p []byte) error
	ReadFrame() ([]byte, error)
	Close() error
}
