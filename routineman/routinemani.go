package routineman

import "context"

type RoutineMan interface {
	Context() context.Context
	Exiting() bool

	StartRoutine(routine func(ctx context.Context), name string)

	TriggerStop()
	StopAndWait()
	Wait()
}
