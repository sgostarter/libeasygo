package statisticsman

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sgostarter/libeasygo/commerr"
	"github.com/sgostarter/libeasygo/helper"
	"github.com/sgostarter/libeasygo/statisticsman/counter"
	"github.com/sgostarter/libeasygo/statisticsman/impl"
	"github.com/sgostarter/libeasygo/statisticsman/inters"
	"github.com/sgostarter/libeasygo/ut"
	"github.com/stretchr/testify/assert"
)

type testDataKey struct {
	Key1 string
	Key2 int
}

func (tsk *testDataKey) Key() string {
	return fmt.Sprintf("%s:%d", tsk.Key1, tsk.Key2)
}

func (tsk *testDataKey) From(s string) error {
	ps := strings.SplitN(s, ":", 2)
	if len(ps) != 2 {
		return commerr.ErrInvalidArgument
	}

	key1 := ps[0]
	key2, err := strconv.Atoi(ps[1])

	if err != nil {
		return err
	}

	tsk.Key1 = key1
	tsk.Key2 = key2

	return nil
}

func TestStatisticsMan(t *testing.T) {
	cfg := ut.SetupUTConfig4Redis(t)
	redisCli, err := helper.NewRedisClient(cfg.RedisDNS)
	assert.Nil(t, err)
	assert.NotNil(t, redisCli)

	r := NewStatisticsReaderEx(impl.NewRedisDataProvider(redisCli), impl.NewTimeSpan(time.Second*5), "")

	tsCounter := counter.NewTimeSpanCounters(impl.NewTimeSpan(time.Second * 5))
	w := NewStatisticsWriterEx(tsCounter, counter.NewAsyncStore(context.Background(), impl.NewRedisCounterStorage(redisCli, nil),
		tsCounter, nil))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		loop := true
		for loop {
			select {
			case <-ctx.Done():
				loop = false

				continue
			case <-time.After(time.Millisecond * 100):
				// nolint:gosec
				n := rand.Intn(2)
				dk := testDataKey{
					Key1: fmt.Sprintf("key%d", n),
					Key2: n,
				}
				w.Inc(&dk)
			}
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		loop := true
		for loop {
			select {
			case <-ctx.Done():
				loop = false

				continue
			case <-time.After(time.Second):
				dk := &testDataKey{}
				m := make(map[string]map[string]int64)
				err = r.Scan4Current(dk, func(timeSpanS string, k inters.DataKey, v int64, err error) error {
					assert.Nil(t, err)
					if _, ok := m[timeSpanS]; !ok {
						m[timeSpanS] = make(map[string]int64)
					}
					m[timeSpanS][k.Key()] = v

					return nil
				})

				for s, i := range m {
					t.Logf("== %s ==\n", s)

					for s2, i2 := range i {
						t.Logf("  %s: %d\n", s2, i2)
					}

					t.Log("== -- ==\n")
				}

				assert.Nil(t, err)
			}
		}
	}()

	wg.Wait()
}

func TestStatisticsMan2(t *testing.T) {
	cfg := ut.SetupUTConfig4Redis(t)
	redisCli, err := helper.NewRedisClient(cfg.RedisDNS)
	assert.Nil(t, err)
	assert.NotNil(t, redisCli)

	r := NewStatisticsReaderEx(impl.NewRedisDataProvider(redisCli), impl.NewTimeSpan(time.Second*5), "")

	dk := &testDataKey{}
	err = r.FlushAndRemoveLastHourData(dk, 10000, func(timeSpanS string, k inters.DataKey, v int64, err error) error {
		t.Log(timeSpanS, k, v, err)

		return nil
	})
	assert.Nil(t, err)
}
