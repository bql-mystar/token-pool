// 原始版本使用的是将redis以key-value形式存入，但这样获取的时候会比较麻烦，因此v2使用list的方式存储token，当需要获取token的使用，弹出一个即可

package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/go-basic/uuid"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"sync"
)

// 定义三个管道用于存储token
var tokenChan = make(chan string, 256)
var wp = sync.WaitGroup{}
var client *redis.Client

// pipeline一次性上传的条数
const execNums = 256

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

// sha1加密工具
func sha1Encryption(content string) string {
	// 对传入内容进行加密
	var h = sha1.New()
	h.Write([]byte(content))
	sha1Content := fmt.Sprintf("%x", h.Sum(nil))
	return sha1Content
}

// 生成token
func generateToken(times int) {
	for i := 0; i < times; i++ {
		// 生成一个uuid
		uuid := uuid.New()
		token := sha1Encryption(uuid)
		// 将token写入channel
		tokenChan <- token
		wp.Add(1)
	}
	close(tokenChan)
}

// 将token写入redis中
func writeToRedis() {
	nums := 0
	// 使用事务pipeline
	pipeline := client.TxPipeline()
	for{
		token, ok := <- tokenChan
		if !ok {
			// 说明通道已经关闭，且里面已经没有值了，直接对将管道中的数据写入redis中即可，也就是flush操作
			// redis管道的写操作
			_, err := pipeline.Exec()
			if err != nil {
				panic(fmt.Sprintf("redis pipeline exec error err is :%s", err.Error()))
			}
			// 同时在waitGroup中减少对应的数据组数
			wp.Add(-1 * nums)
			break
		}
		// 说明读取到内容，将内容通过redis管道的方式写入
		pipeline.LPush(token_list, token)
		nums += 1
		// redis管道的写操作，写入之后对nums的数量+1
		// 如果往通道中加入的数据已经有execNums条的时候，执行redis的flush操作，将其刷入redis中
		if nums >= execNums{
			// 执行redis的flush操作，执行成功后，并将对应的wp的组减少对应的输卵管，并重新将nums置为0
			// flush
			_, err := pipeline.Exec()
			if err != nil {
				panic(fmt.Sprintf("redis pipeline exec error err is :%s", err.Error()))
			}
			wp.Add(-1 * nums)
			nums = 0
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("缺少指定参数")
	}

	// 生成token的数量
	timeStr := os.Args[1]
	times, err := strconv.Atoi(timeStr)

	if err != nil {
		panic("缺少异常")
	}

	// 开启5个工作协程
	for i := 0; i < 10; i++ {
		go writeToRedis()
	}

	// 异步生成token
	generateToken(times)

	wp.Wait()

}

