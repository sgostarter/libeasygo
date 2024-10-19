package fallback

import (
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/l"
)

func NewFallback(policy Policy, logger l.Wrapper) Fallback {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	logger = logger.WithFields(l.StringField(l.ClsKey, "fallbackImpl"))

	if policy == nil {
		logger.Error("no policy")

		return nil
	}

	impl := &fallbackImpl{
		logger: logger,
		policy: policy,
	}

	impl.init()

	return impl
}

type fallbackImpl struct {
	logger l.Wrapper
	policy Policy
}

func (impl *fallbackImpl) init() {

}

func (impl *fallbackImpl) Do(id string, fn func(id string, useFallback bool) error) (useFallback bool, err error) {
	if fn == nil {
		impl.logger.Error("no fn")

		err = commerr.ErrReject

		return
	}

	if !impl.policy.Try(id) {
		err = fn(id, false)
		if err == nil {
			impl.policy.ResultSuccess(id)

			return
		}

		if !impl.policy.ResultError(id, err) {
			return
		}
	}

	useFallback = true

	err = fn(id, true)
	if err == nil {
		impl.policy.ResultFallbackSuccess(id)

		return
	}

	impl.policy.ResultFallbackFailed(id, err)

	return
}
