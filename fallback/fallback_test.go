package fallback

import (
	"testing"
	"time"

	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/l"
)

func TestFallback(t *testing.T) {
	fb := NewFallback(NewDefaultPolicy(nil, &PolicyConfig{
		MaxContinueFailCount:  3,
		TryIntervalOnFailMode: time.Second,
		DataExpiration:        time.Minute,
	}), l.NewConsoleLoggerWrapper())

	var okFlag bool

	fnDo := func(_ string, useFallback bool) error {
		if useFallback {
			t.Log("doFallback")
		} else {
			t.Log("do")
		}

		if okFlag {
			return nil
		}

		return commerr.ErrFailed
	}

	start := time.Now()
	for time.Since(start) < time.Second*5 {
		_, _ = fb.Do("1", fnDo)

		time.Sleep(time.Millisecond * 100)

		if time.Since(start) >= time.Second*3 {
			okFlag = true
		}
	}
}
