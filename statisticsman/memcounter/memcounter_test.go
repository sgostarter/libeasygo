package memcounter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func printMc(t *testing.T, mc *MemCounter) {
	var info string
	for _, d := range mc.ds {
		info += fmt.Sprintf("[%d %d] ", d.tm, d.v)
	}
	t.Log(info)
}

// nolint
func TestMemCounter(t *testing.T) {
	mc := NewMemCounter(time.Second, 6)
	printMc(t, mc)

	mc.Inc()
	printMc(t, mc)

	mc.Inc()
	printMc(t, mc)

	time.Sleep(time.Second)
	mc.Inc()
	printMc(t, mc)

	assert.EqualValues(t, 3, mc.Count())

	time.Sleep(time.Second * 4)
	mc.Inc()
	printMc(t, mc)

	time.Sleep(time.Second * 5)
	mc.Inc()
	printMc(t, mc)

	time.Sleep(time.Second * 6)
	mc.Inc()
	printMc(t, mc)
}

// nolint
func TestMemCounter2(t *testing.T) {
	rand.Seed(time.Now().Unix())
	mc := NewMemCounter(time.Second, 6)

	startAt := time.Now()
	for time.Since(startAt) < time.Minute*5 {
		sleepN := rand.Intn(10)
		time.Sleep(time.Second * time.Duration(sleepN))
		mc.Inc()
		t.Log(sleepN, mc.Count())
	}
}
