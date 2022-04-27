package go_logger

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisWriter struct {
	cli     *redis.Client
	listKey string
	ctx     context.Context
}

func (w *RedisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.ctx, w.listKey, p).Result()
	return int(n), err
}

func (w *RedisWriter) SetListKey(key string) {
	w.listKey = key
}

func NewRedisWriter(ctx context.Context, cli *redis.Client) *RedisWriter {
	return &RedisWriter{
		cli: cli,
		ctx: ctx,
	}
}
