package config

import (
	"fmt"

	"try-on/internal/pkg/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Addr     string
	SqlDir   string
	Redis    Redis
	Postgres Postgres
	Session  Session
	Cors     Cors
}

type Cors struct {
	Domain           string
	AllowCredentials bool
	MaxAge           int
	AllowMethods     []string
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

func NewDynamicConfig(configPath string, onChange func(*Config), onError func(error)) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.BindEnv("postgres.password")

	cfg := Config{}

	viper.OnConfigChange(func(in fsnotify.Event) {
		tmp := Config{}
		err := viper.Unmarshal(&tmp)
		if err != nil {
			if onError != nil {
				onError(err)
			}
			return
		}

		cfg = tmp

		if onChange != nil {
			onChange(&cfg)
		}
	})

	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

var JsonLogFormat = func() string {
	values := []string{"time", "status", "latency", "ip", "method", "path", "error"}

	result := utils.Reduce(values, func(first, second string) string {
		return first + fmt.Sprintf(`"%s": "${%s}",`, second, second)
	})

	return `{"level":"info",` + result[:len(result)-1] + "}\n"
}()

const TimeFormat = "15:04:05 02.01.2006"