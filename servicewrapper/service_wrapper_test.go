package servicewrapper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sgostarter/i/l"
	"github.com/stretchr/testify/assert"
)

type TestCycleJob struct {
	id int
}

func (job *TestCycleJob) OnStart(_ l.Wrapper) {

}

func (job *TestCycleJob) OnFinish(_ l.Wrapper) {

}

func (job *TestCycleJob) DoJob(_ context.Context, _ l.Wrapper) (time.Duration, error) {
	// nolint: forbidigo
	fmt.Printf("[%v]id: %v\n", time.Now(), job.id)

	return time.Second * time.Duration(job.id), nil
}

func TestCycleServiceWrapper(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	sw := NewCycleServiceWrapper(ctx, nil)

	err = sw.Start(&TestCycleJob{id: 1})
	assert.Nil(t, err)

	// err = sw.Start(&TestCycleJob{id: 2})
	// assert.Nil(t, err)

	sw.Wait()
}
