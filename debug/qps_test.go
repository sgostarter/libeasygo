package debug

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestQPS(t *testing.T) {
	qps := New(4)

	var wg sync.WaitGroup

	tEnd := time.Now().Add(time.Second * 8)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for idx := 0; time.Until(tEnd) >= 0; idx++ {
			qps.Inc()

			s := int(time.Until(tEnd).Seconds())
			time.Sleep(time.Millisecond * time.Duration(s))
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		for time.Until(tEnd) >= 0 {
			var s string
			for _, i := range qps.Gets() {
				s += fmt.Sprintf("%d -", i)
			}

			s += "\n"
			t.Log(s)
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
}
