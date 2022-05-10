package configs

import "fmt"

type Redis struct {
	Host     string `json:"host" env:"REDIS_HOST"`
	Port     int    `json:"port" env:"REDIS_PORT"`
	Password string `envconfig:"REDIS_PASSWORD"`
}

// RedisURL return redis connection URL.
func RedisURL() string {
	return fmt.Sprintf("%v:%v", app.Redis.Host, app.Redis.Port)
}

func RedisPassword() string {
	return app.Redis.Password
}
