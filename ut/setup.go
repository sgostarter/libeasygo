package ut

import (
	"testing"

	"github.com/sgostarter/libconfig"
)

func SetupUTConfig() *Config {
	return SetupUTConfigEx("config.yaml", nil)
}

func SetupUTConfigEx(fileName string, configPaths []string) *Config {
	cfg := &Config{}
	var err error
	if len(configPaths) == 0 {
		_, err = libconfig.Load(fileName, cfg)
	} else {
		_, err = libconfig.LoadOnConfigPath(fileName, configPaths, cfg)
	}

	if err != nil {
		return nil
	}

	return cfg
}

func SetupUTConfig4Redis(t *testing.T) *Config {
	return SetupUTConfig4RedisEx("config.yaml", nil, t)
}

func SetupUTConfig4RedisEx(fileName string, configPaths []string, t *testing.T) *Config {
	cfg := SetupUTConfigEx(fileName, configPaths)
	if cfg == nil || cfg.RedisDNS == "" {
		t.SkipNow()

		return nil
	}

	return cfg
}
