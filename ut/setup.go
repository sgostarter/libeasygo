package ut

import (
	"testing"

	"github.com/sgostarter/libconfig"
)

func SetupUTConfig() *Config {
	return SetupUTConfigEx("config.yaml")
}

func SetupUTConfigEx(fileName string) *Config {
	cfg := &Config{}
	_, err := libconfig.Load(fileName, cfg)

	if err != nil {
		return nil
	}

	return cfg
}

func SetupUTConfig4Redis(t *testing.T) *Config {
	return SetupUTConfig4RedisEx("config.yaml", t)
}

func SetupUTConfig4RedisEx(fileName string, t *testing.T) *Config {
	cfg := SetupUTConfigEx(fileName)
	if cfg == nil || cfg.RedisDNS == "" {
		t.SkipNow()

		return nil
	}

	return cfg
}
