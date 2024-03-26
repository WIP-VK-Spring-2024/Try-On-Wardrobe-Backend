package config

import (
	"fmt"
	"time"

	"try-on/internal/pkg/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Addr       string
	Static     Static
	Centrifugo Centrifugo
	Postgres   Postgres
	Session    Session
	Cors       Cors
	Sql        Sql
	S3         S3
	Rabbit     Rabbit
}

type Static struct {
	HttpApi  HttpApi
	Type     string
	Dir      string
	Clothes  string
	Cut      string
	FullBody string
	TryOn    string
	S3       S3
}

type HttpApi struct {
	Endpoint    string
	Token       string
	TokenHeader string
	UploadUrl   string
	DeleteUrl   string
}

type Centrifugo struct {
	Url               string
	TryOnChannel      string
	ProcessingChannel string
}

type Cors struct {
	Domain           string
	AllowCredentials bool
	MaxAge           int
	AllowMethods     []string
}

type Rabbit struct {
	Host     string
	Port     int
	User     string
	Password string
	TryOn    RabbitQueue
	Process  RabbitQueue
}

type RabbitQueue struct {
	Request  string
	Response string
}

func (cfg *Rabbit) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
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

func (cfg *Postgres) PoolDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable pool_max_conns=%d",
		cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port, cfg.MaxConn)
}

func (cfg *Postgres) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port)
}

type Sql struct {
	Dir string
}

type Session struct {
	TokenName string
	MaxAge    int
	Secret    string
}

type S3 struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func NewDynamicConfig(configPath string, onChange func(*Config), onError func(error)) (*Config, error) {
	viper.SetConfigFile(configPath)

	viper.BindEnv("postgres.host")
	viper.BindEnv("postgres.port")
	viper.BindEnv("postgres.password")
	viper.BindEnv("session.secret")
	viper.BindEnv("static.s3.accessKey")
	viper.BindEnv("static.s3.secretKey")
	viper.BindEnv("rabbit.password")
	viper.BindEnv("static.httpapi.token")
	viper.BindEnv("rabbit.host")
	viper.BindEnv("static.httpapi.endpoint")

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
