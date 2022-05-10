package routineman

import (
	"context"
	"time"
)

type RoutineMan interface {
	Context() context.Context
	Exiting() bool

	StartRoutine(routine func(ctx context.Context, exiting func() bool), name string)

	TriggerStop()
	StopAndWait()
	Wait()

	DoWithTimeout(label string, do func(), d time.Duration)
}

type DebugRoutineManTimeoutObserver func(msg string)

type DebugRoutineMan interface {
	RoutineMan

	SetExitTimeoutObserver(ob DebugRoutineManTimeoutObserver)
}
