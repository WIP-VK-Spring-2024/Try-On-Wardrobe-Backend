package config

import (
	"fmt"
	"time"

	"try-on/internal/pkg/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Addr     string
	SqlDir   string
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

type Postgres struct {
	DB          string
	User        string
	Password    string
	Host        string
	Port        string
	MaxConn     int
	InitTimeout time.Duration
}

func (cfg *Postgres) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port)
}

type Session struct {
	TokenName string
	MaxAge    int
	Secret    string
}

func NewDynamicConfig(configPath string, onChange func(*Config), onError func(error)) (*Config, error) {
	viper.SetConfigFile(configPath)

	viper.BindEnv("postgres.host")
	viper.BindEnv("postgres.port")
	viper.BindEnv("postgres.password")
	viper.BindEnv("session.secret")

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
	values := []string{"time", "status", "latency", "ip", "method", "path"}

	result := utils.Reduce(values, func(first, second string) string {
		return first + fmt.Sprintf(`"%s":"${%s}",`, second, second)
	})

	return `{"level":"info",` + result[:len(result)-1] + "}\n"
}()

const TimeFormat = "15:04:05 02.01.2006"
