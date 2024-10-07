package flowlimit

type Limiter interface {
	Allow(key string) bool

	Close()
	Wait()
}

type Hystrix interface {
	Do(key string, d func() error, f func(err error)) error
	Go(key string, d func() error, f func(err error)) error
}
