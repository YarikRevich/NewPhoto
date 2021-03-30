package caching

import (
	"context"
	"math"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var(
	RedisInstanse = New()
	ctx = context.Background()
)

type Redis struct{
	db *redis.Client
}

func (r *Redis) Connect(){
	r.db = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "hello",
		DB: 0,
		Username: "default",
	})
	if err := r.db.Ping(ctx).Err(); err != nil{
		log.Fatalln(err)
	}
}

func (r *Redis) Set(key, rc ,args string){
	if err := r.db.HSet(ctx, key, []string{"commands", rc}).Err(); err != nil{
		log.Fatalln(err)
	}
	if err := r.db.HSet(ctx, key, []string{"args", args}).Err(); err != nil{
		log.Fatalln(err)
	}
	expr := time.Now().Add(time.Second * 10).Format("2006-01-02 15:04:05-07:00")
	if err := r.db.HSet(ctx, key, []string{"expr", expr}).Err(); err != nil{
		log.Fatalln(err)
	}
}

func (r *Redis) Get(key string)(string, string, string){
	command, err := r.db.HGet(ctx, key, "commands").Result()
	if err != nil{
		command = ""
	}
	args, err := r.db.HGet(ctx,  key, "args").Result()
	if err != nil{
		args = ""
	}
	expr, err := r.db.HGet(ctx, key, "expr").Result()
	if err != nil{
		expr = "2006-01-02 15:04:05-07:00"
	}
	return command, args, expr
}

func (r *Redis) IsCached(key, expectedcommand string)(string, bool){
	command, args, expr := r.Get(key)
	ttl, err := time.Parse("2006-01-02 15:04:05-07:00", expr)
	if err != nil{
		log.Fatalln(err)
	}
	cexpr := ttl.Sub(time.Now())
	if len(command) != 0 &&  (command == expectedcommand && !math.Signbit(cexpr.Minutes())){
		return args, true
	}
	r.db.Del(ctx, key)
	return "nil", false
}

func New()*Redis{
	r := new(Redis)
	r.Connect()
	return r
}
