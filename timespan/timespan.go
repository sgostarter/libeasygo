package timespan

import "time"

const (
	humanLayout = "20060102150405"
)

type TimeSpan struct {
	durationSeconds int64
	location        *time.Location
	locationOffset  int
}

func NewTimeSpan(duration time.Duration) *TimeSpan {
	return NewTimeSpanEx(duration, nil)
}

func NewTimeSpanEx(duration time.Duration, location *time.Location) *TimeSpan {
	if location == nil {
		location = time.Local
	}

	_, locationOffset := time.Now().In(location).Zone()

	if duration < time.Second {
		panic("invalid duration")
	}

	return &TimeSpan{
		durationSeconds: int64(duration / time.Second),
		location:        location,
		locationOffset:  locationOffset,
	}
}

func (ts *TimeSpan) GetLabel(t time.Time) string {
	tUnix := t.Unix() + int64(ts.locationOffset)

	return time.Unix(tUnix-tUnix%ts.durationSeconds-int64(ts.locationOffset), 0).Format(humanLayout)
}

func (ts *TimeSpan) GetCurrentLabel() string {
	return ts.GetLabel(time.Now())
}

func (ts *TimeSpan) DiffTime(timeStart, timeFinish time.Time) int {
	return int((timeFinish.Unix() - timeStart.Unix()) / ts.durationSeconds)
}

func (ts *TimeSpan) DiffLabel(labelStart, labelFinish string) (n int, err error) {
	startTime, err := ts.Label2Time(labelStart)
	if err != nil {
		return
	}

	finishTime, err := ts.Label2Time(labelFinish)

	if err != nil {
		return
	}

	n = ts.DiffTime(startTime, finishTime)

	return
}

func (ts *TimeSpan) Label2Time(l string) (t time.Time, err error) {
	return time.ParseInLocation(humanLayout, l, time.Local)
}

func (ts *TimeSpan) LabelString2Int(l string) (n int64, err error) {
	t, err := ts.Label2Time(l)
	if err != nil {
		return
	}

	n = t.Unix()

	return
}

func (ts *TimeSpan) LabelInt2String(n int64) string {
	return ts.GetLabel(time.Unix(n, 0))
}
