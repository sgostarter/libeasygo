package ut

type Config struct {
	RedisDNS string `yaml:"redis_dns" envconfig:"REDIS_DSN"`
}
