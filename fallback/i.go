package fallback

type Policy interface {
	Try(id string) (useFallback bool)
	ResultSuccess(id string)
	ResultError(id string, err error) (useFallback bool)
	ResultFallbackSuccess(id string)
	ResultFallbackFailed(id string, err error)
}

type Fallback interface {
	Do(id string, fb func(id string, useFallback bool) error) (useFallback bool, err error)
}
