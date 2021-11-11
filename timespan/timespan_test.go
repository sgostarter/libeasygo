package timespan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeSpan(t *testing.T) {
	ts := NewTimeSpan(time.Second)
	t.Log(ts.GetCurrentLabel())
}

func TestTimeSpan1Hour(t *testing.T) {
	ts := NewTimeSpan(time.Hour)

	st, err := time.ParseInLocation(humanLayout, "20211103105501", time.Local)
	assert.Nil(t, err)

	sl := ts.GetLabel(st)

	nt, err := time.ParseInLocation(humanLayout, "20211103115501", time.Local)
	assert.Nil(t, err)

	nl := ts.GetLabel(nt)

	nt2, err := time.ParseInLocation(humanLayout, "20211103105959", time.Local)
	assert.Nil(t, err)

	nl2 := ts.GetLabel(nt2)

	nt3, err := time.ParseInLocation(humanLayout, "20211103110000", time.Local)
	assert.Nil(t, err)

	nl3 := ts.GetLabel(nt3)

	t.Log(sl)
	t.Log(nl)

	n, err := ts.DiffLabel(sl, nl)
	assert.Nil(t, err)
	assert.Equal(t, n, 1)

	n, err = ts.DiffLabel(sl, nl2)
	assert.Nil(t, err)
	assert.Equal(t, n, 0)

	n, err = ts.DiffLabel(sl, nl3)
	assert.Nil(t, err)
	assert.Equal(t, n, 1)
}

func TestTimeSpan5Minute(t *testing.T) {
	ts := NewTimeSpan(time.Minute * 5)

	st, err := time.ParseInLocation(humanLayout, "20211103101521", time.Local)
	assert.Nil(t, err)

	sl := ts.GetLabel(st)

	nt, err := time.ParseInLocation(humanLayout, "20211103102021", time.Local)
	assert.Nil(t, err)

	nl := ts.GetLabel(nt)

	nt2, err := time.ParseInLocation(humanLayout, "20211103101959", time.Local)
	assert.Nil(t, err)

	nl2 := ts.GetLabel(nt2)

	nt3, err := time.ParseInLocation(humanLayout, "20211103102001", time.Local)
	assert.Nil(t, err)

	nl3 := ts.GetLabel(nt3)

	t.Log(sl)
	t.Log(nl)
	t.Log(nl2)
	t.Log(nl3)

	n, err := ts.DiffLabel(sl, nl)
	assert.Nil(t, err)
	assert.Equal(t, n, 1)

	n, err = ts.DiffLabel(sl, nl2)
	assert.Nil(t, err)
	assert.Equal(t, n, 0)

	n, err = ts.DiffLabel(sl, nl3)
	assert.Nil(t, err)
	assert.Equal(t, n, 1)
}

func TestTimeSpan4Day(t *testing.T) {
	ts := NewTimeSpan(time.Hour * 24)

	timeNow := time.Now()

	days, err := ts.DiffLabel(ts.GetLabel(timeNow), ts.GetLabel(timeNow.Add(time.Hour)))
	assert.Nil(t, err)
	assert.Equal(t, days, 0)

	days, err = ts.DiffLabel(ts.GetLabel(timeNow), ts.GetLabel(timeNow.Add(24*time.Hour)))
	assert.Nil(t, err)
	assert.Equal(t, days, 1)

	tm1, err := time.Parse("2006-01-02 15:04:05", "2021-11-01 15:11:52")
	assert.Nil(t, err)

	tm2, err := time.Parse("2006-01-02 15:04:05", "2021-11-01 23:59:59")
	assert.Nil(t, err)

	tm3, err := time.Parse("2006-01-02 15:04:05", "2021-11-02 00:00:00")
	assert.Nil(t, err)

	tm4, err := time.Parse("2006-01-02 15:04:05", "2021-12-01 15:11:52")
	assert.Nil(t, err)

	days, err = ts.DiffLabel(ts.GetLabel(tm1), ts.GetLabel(tm2))
	assert.Nil(t, err)
	assert.Equal(t, days, 0)

	days, err = ts.DiffLabel(ts.GetLabel(tm1), ts.GetLabel(tm3))
	assert.Nil(t, err)
	assert.Equal(t, days, 1)

	days, err = ts.DiffLabel(ts.GetLabel(tm1), ts.GetLabel(tm4))
	assert.Nil(t, err)
	assert.Equal(t, days, 30)
}
