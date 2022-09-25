package debug

import (
	"os"
	"testing"
)

func TestRunProfileProxy(t *testing.T) {
	t.SkipNow()

	_ = RunProfileProxy(":10000", os.Getenv("token"))
}
