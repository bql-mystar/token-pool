package main

// 从token池中返回一个未使用的token，并将token标记为已使用

import "github.com/go-redis/redis"

var client *redis.Client

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

func getTokenFromPool(){}


func main(){

}
