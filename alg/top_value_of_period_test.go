package alg

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/sgostarter/libeasygo/i"
)

// nolint
func TestTopValueOfPeriod(t *testing.T) {
	c := NewTopValueOfPeriod[int64, i.RWLock](&sync.RWMutex{}, time.Second*6, &MinCheck[int64]{})

	for idx := 0; idx < 10; idx++ {
		v := rand.Int31n(30)
		c.Set(int64(v))
		topV, exists := c.Get()
		t.Log("new add is:", v, "top is:", topV, exists)
		time.Sleep(time.Second)
	}
}
