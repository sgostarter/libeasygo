package memcounter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	t.SkipNow()

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
	t.SkipNow()

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

// nolint
func TestMemCount3(t *testing.T) {
	t.SkipNow()

	rand.Seed(time.Now().Unix())
	mc := NewMemCounter(time.Second, 6)

	counts := make(map[int64]int)

	var compareCount int

	for idx := 0; idx < 1000; idx++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))

		mc.Inc()
		counts[time.Now().Unix()]++

		if idx%7 != 0 {
			continue
		}

		firstTm := time.Now().Unix() - 5

		count1 := mc.Count()

		var count2 int

		for tm, c := range counts {
			if tm < firstTm {
				delete(counts, tm)
			} else {
				count2 += c
			}
		}

		assert.EqualValues(t, count2, count1)

		compareCount++

		if compareCount%10 == 0 {
			t.Log("Compare: ", compareCount)
		}
	}
}
