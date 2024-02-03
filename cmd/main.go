package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/christianferraz/ratelimiter/configs"
	"github.com/christianferraz/ratelimiter/internal/entity"
	"github.com/christianferraz/ratelimiter/limiter"
	"github.com/christianferraz/ratelimiter/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := configs.LoadConfig(".")
	ctx := context.Background()
	if err != nil {
		panic(err)
	}
	rd_client := redis.NewClient(&redis.Options{
		Addr:     config.RedisSrc,
		Password: config.RedisPass,
		DB:       0,
	})
	defer rd_client.Close()
	pong, err := rd_client.Ping(ctx).Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Conexão ao Redis estabelecida:", pong)
	rds := entity.NewRedisStorage(rd_client)
	rl := limiter.NewRateLimiter(config, rds)
	http.HandleFunc("/", middleware.CountMiddleware(handler, rl))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
