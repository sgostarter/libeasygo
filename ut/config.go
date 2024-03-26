package ut

import (
	"testing"
)

type Config struct {
	RedisDNS string         `yaml:"redis_dns" envconfig:"REDIS_DSN"`
	Envs     map[string]any `json:"envs" yaml:"envs" envconfig:"ENVS"`
}

func (cfg *Config) GetEnv(t *testing.T, key string) any {
	v, ok := cfg.Envs[key]
	if !ok {
		t.SkipNow()
	}

	return v
}
