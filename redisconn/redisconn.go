package redisconn

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)



var inited bool = false
var RedisClient *redis.Pool
const rediscfg string = "127.0.0.1:6379"

func InitPool() {
	if (!inited) {
		RedisClient = &redis.Pool{
			//连接方法
			Dial:            func() (redis.Conn,error) {
				c,err := redis.Dial("tcp", rediscfg)
				if err != nil {
					return nil,err
				}
				c.Do("SELECT",0)
				return c,nil
			},
			//DialContext:     nil,
			//TestOnBorrow:    nil,
			//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
			MaxIdle:         1,
			//最大的激活连接数，表示同时最多有N个连接
			MaxActive:       10,
			//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
			IdleTimeout:     180 * time.Second,
			//Wait:            false,
			//MaxConnLifetime: 0,
			}
	}
}

func test() {
	c1 := RedisClient.Get()
	defer c1.Close()

	c1.Do("SET", "xhc", "123")

	r, err := c1.Do("GET","xhc")
	if err == nil {
		fmt.Println(string(r.([]byte)))
	}
}

func SetProfile(username string, profile string) {

	c1 := RedisClient.Get()
	defer c1.Close()
	_, err := c1.Do("SET", "profile" + username, profile)

	if err != nil {
		fmt.Println("setProfile set error ", err)
	}
}

func GetProfile(username string) (string, error) {
	c1 := RedisClient.Get();
	defer c1.Close()

	res, err := c1.Do("GET", "profile" + username)
	if  err != nil {
		fmt.Println("set err = ", err)
		return "", err
	}

	return string(res.([]byte)), nil
}

func SetToken(username string, token string) {

	c1 := RedisClient.Get()
	defer c1.Close()
	_, err := c1.Do("SET", "token" + username, token)

	if err != nil {
		fmt.Println("SetToken set error ", err)
	}
}

func GetToken(username string) string {
	c1 := RedisClient.Get();
	defer c1.Close()

	res, err := c1.Do("GET", "token" + username)
	if  err != nil || res == nil {
		fmt.Println("GetToken err = ", err)
		return ""
	}
	return string(res.([]byte))
}