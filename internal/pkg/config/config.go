package config

import (
	"fmt"

	"try-on/internal/pkg/utils"
)

type Config struct {
	Port string
	Redis
	Postgres
	Session
	Cors
}

type Cors struct {
	Domain           string
	AllowCredentials bool
	MaxAge           int
}

type Redis struct {
	Addr    string
	MaxConn int
}

type Postgres struct {
	DB       string
	User     string
	Password string
	Host     string
	Port     string
	MaxConn  int
}

func (cfg *Postgres) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port)
}

type Session struct {
	KeyNamespace string
	CookieName   string
	MaxAge       int
}

var JsonLogFormat = func() string {
	values := []string{"time", "status", "latency", "ip", "method", "path", "error"}

	result := utils.Reduce(values, func(first, second string) string {
		return first + fmt.Sprintf(`%s: "${%s}", `, second, second)
	})

	return "{" + result[:len(result)-2] + "}\n"
}()
