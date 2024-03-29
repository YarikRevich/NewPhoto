package caching

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/YarikRevich/NewPhoto/log"
	"github.com/go-redis/redis/v8"
)

var (
	RedisInstanse = New()
	ctx           = context.Background()
)

type Redis struct {
	db *redis.Client
}

func (r *Redis) Connect() {
	addr, ok := os.LookupEnv("redisAddr")
	if !ok {
		log.Logger.UsingErrorLogFile().CFatalln("ChacheInit", "redisAddr is not written in credentials.sh file")
	}

	password, ok := os.LookupEnv("redisPassword")
	if !ok {
		log.Logger.UsingErrorLogFile().CFatalln("ChacheInit", "redisPassword is not written in credentials.sh file")
	}

	username, ok := os.LookupEnv("redisUsername")
	if !ok {
		log.Logger.UsingErrorLogFile().CFatalln("ChacheInit", "redisUsername is not written in credentials.sh file")
	}

	r.db = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", addr),
		Password: password,
		DB:       0,
		Username: username,
	})
	if err := r.db.Ping(ctx).Err(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("ChacheInit", err)
	}
}

func (r *Redis) Set(key, rc, args string) {
	if err := r.db.HSet(ctx, key, []string{"commands", rc}).Err(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("ChacheSet", err)
	}
	if err := r.db.HSet(ctx, key, []string{"args", args}).Err(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("ChacheSet", err)
	}
	expr := time.Now().Add(time.Second * 10).Format("2006-01-02 15:04:05-07:00")
	if err := r.db.HSet(ctx, key, []string{"expr", expr}).Err(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("ChacheSet", err)
	}
}

func (r *Redis) Get(key string) (string, string, string) {
	command, err := r.db.HGet(ctx, key, "commands").Result()
	if err != nil {
		command = ""
	}
	args, err := r.db.HGet(ctx, key, "args").Result()
	if err != nil {
		args = ""
	}
	expr, err := r.db.HGet(ctx, key, "expr").Result()
	if err != nil {
		expr = "2006-01-02 15:04:05-07:00"
	}
	return command, args, expr
}

func (r *Redis) IsCached(key, expectedcommand string) (string, bool) {
	command, args, expr := r.Get(key)
	ttl, err := time.Parse("2006-01-02 15:04:05-07:00", expr)
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("IsCached", err)
	}

	if len(command) != 0 && (command == expectedcommand && !math.Signbit(time.Until(ttl).Minutes())) {
		return args, true
	}
	r.db.Del(ctx, key)
	return "nil", false
}

func (r *Redis) Clean(key string) {
	if c := r.db.Del(ctx, key); c.Err() != nil {
		log.Logger.CFatalln("CleanCache", c.Err())
	}
}

func New() *Redis {
	r := new(Redis)
	r.Connect()
	return r
}
