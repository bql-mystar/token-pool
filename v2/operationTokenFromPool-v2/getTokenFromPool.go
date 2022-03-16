package main

import "fmt"

// 从token池中获取token
func getTokenFromPool() (string, error){
	token := client.RPop(token_list)
	if token.Err() != nil {
		return "", fmt.Errorf("get token failed")
	}
	// 将token返回，并将对应的token以字符串的形式设置在redis中，如果后期可以通过get方法获取对应的key的值，说明该token已经被使用
	client.Set(token.Val(), "1", 0)
	return token.Val(), nil
}
