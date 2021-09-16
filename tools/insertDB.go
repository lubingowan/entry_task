package main

import (
	"entry_task/mysqlconn"
	"entry_task/token"
	"fmt"
	"strconv"
)

func main() {
	username := "test"
	nickname := "Rubin"
	password := "password"
	
	for i := 0; i < 5000; i++ {
		istr := strconv.Itoa(i)
		p := &mysqlconn.Profile{
			Username: username + istr,
			Nickname: nickname + istr,
			Password: []byte(token.Sha1(password + istr)),
		}

		mysqlconn.InsertProfile(p)
		fmt.Println(p)
	}
}