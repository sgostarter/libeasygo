package debug

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	rt := NewRT(context.Background(), time.Minute)

	wg := sync.WaitGroup{}

	fnDo := func() {
		defer wg.Done()

		probe := rt.Enter()
		defer probe.Leave()

		time.Sleep(time.Millisecond * 100)
	}

	wg.Add(10)

	for idx := 0; idx < 10; idx++ {
		go fnDo()
	}

	wg.Wait()

	time.Sleep(time.Second)

	avg := rt.GetAVG()
	assert.True(t, avg >= time.Millisecond*100 && avg < time.Millisecond*110)

	rt.DestroyAndWait()
}
