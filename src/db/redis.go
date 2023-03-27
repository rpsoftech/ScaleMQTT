package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisConn *redis.Client

// func A() {

// 	redisConn := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "",
// 		DB:       0,
// 	})
// }
