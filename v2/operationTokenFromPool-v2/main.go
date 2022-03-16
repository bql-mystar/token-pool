package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

var client *redis.Client

// token的相关操作

// token池队列名称
const token_list = "token-pool"

// 初始化redis客户端
func init() {
	client = redis.NewClient(&redis.Options{
		//Addr: "127.0.0.1:6379",
		Addr: "120.25.220.19:6379",
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic("redis connect error")
	}
}



func main() {
	fmt.Println(getTokenFromPool())
}

