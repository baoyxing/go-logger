package main

import (
	"context"
	go_logger "github.com/baoyxing/go-logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	writer := go_logger.NewRedisWriter(context.Background(), cli)
	writer.SetListKey("test")
	logger := go_logger.InitLogger(nil, zapcore.InfoLevel)
	logger.Error("test logger info", zap.String("hello", "logger"))
}
