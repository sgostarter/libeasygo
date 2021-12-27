package impl

import (
	"time"

	"github.com/sgostarter/libeasygo/statisticsman/counter"
)

func NewHourTimeSpan() counter.TimeSpan {
	return NewTimeSpan(time.Hour)
}
