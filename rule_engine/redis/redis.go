package redis

import (
	"context"
	_redis "github.com/go-redis/redis"
	"time"
)

var (
	client *_redis.Client
)

const queryTimeout = 3 * time.Second

func GetConn(ctx context.Context) *_redis.Client {
	return client.WithContext(ctx)
}

func Init(cfg *Config) error {
	cfg.FillWithDefaults()

	client = _redis.NewClient(&_redis.Options{
		Network:      cfg.Network,
		Addr:         cfg.Addr,
		Password:     cfg.Passwd,
		DB:           cfg.DB,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
	})

	return client.Ping().Err()
}

func Publish(key string, val []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := GetConn(ctx)
	return conn.Publish(key, val).Err()
}
