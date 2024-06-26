package config

import (
	"fmt"
	"time"

	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Addr           string
	WeatherApiKey  string
	Static         Static
	Centrifugo     Centrifugo
	Postgres       Postgres
	Session        Session
	Cors           Cors
	Sql            Sql
	S3             S3
	Rabbit         Rabbit
	Redis          Redis
	Classification Classification
	ModelsHealth   ModelsHealth
}

type Static struct {
	HttpApi         HttpApi
	Type            string
	Dir             string
	Clothes         string
	Cut             string
	FullBody        string
	TryOn           string
	Outfits         string
	Avatars         string
	S3              S3
	DefaultImgPaths DefaultImgPaths
}

type DefaultImgPaths map[domain.Gender]string

type Classification struct {
	Threshold float32
}

type HttpApi struct {
	Endpoint    string
	Token       string
	TokenHeader string
	UploadUrl   string
	DeleteUrl   string
	GetUrl      string
}

type Redis struct {
	Host string
	Port int
}

func (r Redis) DSN() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type Centrifugo struct {
	Url               string
	TryOnChannel      string
	ProcessingChannel string
	OutfitGenChannel  string
}

type ModelsHealth struct {
	Token       string
	Endpoint    string
	TryOn       string
	Cut         string
	OutfitGen   string
	Recsys      string
	TokenHeader string
}

type Cors struct {
	Domain           string
	AllowCredentials bool
	MaxAge           int
	AllowMethods     []string
}

type Rabbit struct {
	Host      string
	Port      int
	User      string
	Password  string
	TryOn     RabbitQueue
	Process   RabbitQueue
	OutfitGen RabbitQueue
	Recsys    RabbitQueue
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

var envBoundConfigValues = []string{
	"postgres.host",
	"postgres.port",
	"postgres.password",
	"session.secret",
	"static.s3.accessKey",
	"static.s3.secretKey",
	"rabbit.password",
	"static.httpapi.token",
	"rabbit.host",
	"static.httpapi.endpoint",
	"modelsHealth.token",
}

func NewDynamicConfig(configPath string, onChange func(*Config), onError func(error)) (*Config, error) {
	viper.SetConfigFile(configPath)

	for _, value := range envBoundConfigValues {
		viper.BindEnv(value)
	}
	viper.BindEnv("weatherapikey", "WEATHER_API_KEY")

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
