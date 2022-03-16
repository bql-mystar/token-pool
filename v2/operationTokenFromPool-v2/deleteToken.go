package main

import "fmt"

func deleteToken(token string)error{
	res := client.Del(token)
	if res.Err() != nil {
		return fmt.Errorf("delete token error")
	}
	return nil
}
