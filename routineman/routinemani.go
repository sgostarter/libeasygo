package routineman

import "context"

type RoutineMan interface {
	Context() context.Context
	Exiting() bool

	StartRoutine(routine func(ctx context.Context, exiting func() bool), name string)

	TriggerStop()
	StopAndWait()
	Wait()
}

type DebugRoutineManTimeoutObserver func(msg string)

type DebugRoutineMan interface {
	RoutineMan

	SetExitTimeoutObserver(ob DebugRoutineManTimeoutObserver)
}
