package cache

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	DB       int    `yaml:"database"`
}
