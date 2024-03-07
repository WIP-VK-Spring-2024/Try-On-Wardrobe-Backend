package redis

import (
	"encoding/json"

	"try-on/internal/pkg/api_errors"
	"try-on/internal/pkg/domain"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Namespace     string
	MaxConn       int
	ExpireSeconds int64
	RedisAddr     string
}

func NewRedisSessionStorage(cfg Config) domain.SessionRepository {
	return &RedisSessionStorage{
		pool: &redis.Pool{
			MaxActive: cfg.MaxConn,
			Dial:      func() (redis.Conn, error) { return redis.Dial("tcp", cfg.RedisAddr) },
		},
		namespace:     cfg.Namespace,
		expireSeconds: cfg.ExpireSeconds,
	}
}

type RedisSessionStorage struct {
	pool          *redis.Pool
	namespace     string
	expireSeconds int64
}

func (repo *RedisSessionStorage) Put(session domain.Session) error {
	conn := repo.pool.Get()
	defer conn.Close()

	return repo.add(conn, session, repo.expireSeconds)
}

func (repo *RedisSessionStorage) Get(key string) (*domain.Session, error) {
	var bytes []byte

	conn := repo.pool.Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("GET", repo.getKey(key)))
	if err != nil {
		if err == redis.ErrNil {
			err = api_errors.ErrNotFound
		}
		return nil, err
	}

	result := &domain.Session{}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *RedisSessionStorage) Delete(key string) error {
	conn := repo.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", repo.getKey(key))
	return err
}

func (repo *RedisSessionStorage) getKey(key string) string {
	return repo.namespace + ":" + key
}

func (repo *RedisSessionStorage) add(conn redis.Conn, session domain.Session, expireSeconds int64) error {
	bytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", repo.getKey(session.ID), bytes, "EX", expireSeconds)
	return err
}
