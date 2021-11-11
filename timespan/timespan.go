package timespan

import "time"

const (
	humanLayout = "20060102150405"
)

type TimeSpan struct {
	durationSeconds int64
}

func NewTimeSpan(duration time.Duration) *TimeSpan {
	return &TimeSpan{
		durationSeconds: int64(duration / time.Second),
	}
}

func (ts *TimeSpan) GetLabel(t time.Time) string {
	return time.Unix(t.Unix()-t.Unix()%ts.durationSeconds, 0).Format(humanLayout)
}

func (ts *TimeSpan) GetCurrentLabel() string {
	return ts.GetLabel(time.Now())
}

func (ts *TimeSpan) DiffTime(timeStart, timeFinish time.Time) int {
	return int((timeFinish.Unix() - timeStart.Unix()) / ts.durationSeconds)
}

func (ts *TimeSpan) DiffLabel(labelStart, labelFinish string) (n int, err error) {
	startTime, err := ts.label2Time(labelStart)
	if err != nil {
		return
	}

	finishTime, err := ts.label2Time(labelFinish)

	if err != nil {
		return
	}

	n = ts.DiffTime(startTime, finishTime)

	return
}

func (ts *TimeSpan) label2Time(l string) (t time.Time, err error) {
	return time.ParseInLocation(humanLayout, l, time.Local)
}
