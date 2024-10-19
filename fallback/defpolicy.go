package fallback

import (
	"fmt"
	"time"

	"github.com/godruoyi/go-snowflake"
	"github.com/patrickmn/go-cache"
	"go.uber.org/atomic"
)

type PolicyConfig struct {
	MaxContinueFailCount  int32
	TryIntervalOnFailMode time.Duration
	DataExpiration        time.Duration
}

func (cfg *PolicyConfig) Valid() {
	if cfg.MaxContinueFailCount <= 0 {
		cfg.MaxContinueFailCount = 1
	}

	if cfg.TryIntervalOnFailMode <= 0 {
		cfg.TryIntervalOnFailMode = time.Minute
	}

	if cfg.DataExpiration <= 0 {
		cfg.DataExpiration = time.Hour
	}
}

func NewDefaultPolicy(dCache *cache.Cache, cfg *PolicyConfig) Policy {
	if cfg == nil {
		cfg = &PolicyConfig{
			MaxContinueFailCount:  2,
			TryIntervalOnFailMode: time.Minute * 2,
			DataExpiration:        time.Hour,
		}
	}

	cfg.Valid()

	if dCache == nil {
		dCache = cache.New(cfg.DataExpiration, cfg.DataExpiration)
	}

	id := snowflake.ID()

	return &defPolicyImpl{
		id:      id,
		keyBase: fmt.Sprintf("default-fallback-policy:%d", id),
		dCache:  dCache,
		cfg:     cfg,
	}
}

type defPolicyImpl struct {
	id      uint64
	keyBase string

	dCache *cache.Cache
	cfg    *PolicyConfig
}

type cacheData struct {
	lastCheckAt       atomic.Time
	continueFailCount atomic.Int32
}

func (impl *defPolicyImpl) cacheKey(id string) string {
	return impl.keyBase + id
}

func (impl *defPolicyImpl) Try(id string) (useFallback bool) {
	i, ok := impl.dCache.Get(impl.cacheKey(id))
	if !ok {
		return
	}

	cd, ok := i.(*cacheData)
	if !ok {
		return
	}

	if cd.continueFailCount.Load() < impl.cfg.MaxContinueFailCount {
		return
	}

	if time.Since(cd.lastCheckAt.Load()) > impl.cfg.TryIntervalOnFailMode {
		cd.lastCheckAt.Store(time.Now())

		return
	}

	useFallback = true

	return
}

func (impl *defPolicyImpl) ResultSuccess(id string) {
	impl.dCache.Delete(impl.cacheKey(id))
}

func (impl *defPolicyImpl) ResultError(id string, _ error) (useFallback bool) {
	cacheKey := impl.cacheKey(id)

	var cd *cacheData

	i, ok := impl.dCache.Get(cacheKey)
	if ok {
		cd, _ = i.(*cacheData)
	}

	if cd == nil {
		cd = &cacheData{}

		impl.dCache.Set(cacheKey, cd, impl.cfg.DataExpiration)
	}

	cd.continueFailCount.Add(1)
	cd.lastCheckAt.Store(time.Now())

	useFallback = true

	return
}

func (impl *defPolicyImpl) ResultFallbackSuccess(_ string) {

}

func (impl *defPolicyImpl) ResultFallbackFailed(_ string, _ error) {

}
