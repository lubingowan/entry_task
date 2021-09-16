package main

import (
	// "bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"entry_task/mysqlconn"
	"entry_task/protocol"
	"entry_task/redisconn"
)

func peekhead(buff []byte) (uri int32, length int32, e error) {
	decbuffer := bytes.NewBuffer(buff)

	err := binary.Read(decbuffer, binary.LittleEndian, &uri)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(decbuffer, binary.LittleEndian, &length)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("peekhead, uri ", uri, " len ", length, " err ", err)

	return uri, length, err
}

func doRead(conn net.Conn) (uri int32, buff []byte) {
	var length int32 = 0
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	uri, length, err = peekhead(buf)
	if err != nil {
		fmt.Println("after peek head, ", err)
		return 0, nil
	}

	fmt.Println("uri: ", uri, " length: ", length, "  n: ", n, " err: ", err)
	fmt.Println(buf)

	if n < int(length + 8) {
		return 0, nil
	}

	return uri, buf[8:8+length]
} 

func onAuthInfo(authinfo protocol.AuthInfo, conn net.Conn) {
	password := redisconn.GetPassword(authinfo.Username)

	fmt.Println("query redis, ", password)

	authres := protocol.AuthResult{
		Username: authinfo.Username, 
		Result: false, 
		Reqid: authinfo.Reqid }

	if len(password) == 0 {
		p, err := mysqlconn.QueryProfile(authinfo.Username)
		if err != nil {
			fmt.Println("onAuthInfo query db failed, ", err)
			authres.Result = false
		} else {
			fmt.Println(p, " authinfo.password ", authinfo.Password, p.Password == authinfo.Password, " ", p.Password, "-", authinfo.Password)
			if p.Password == authinfo.Password {
				authres.Result = true
				authres.Nickname = p.Nickname
				authres.Picture = p.Picture
				defer storeProfile(&p)
			}
		}
	} else {
		if password == authinfo.Password {
			authres.Result = true

			p, err := redisconn.GetProfile(authinfo.Username)
			if err != nil {
				fmt.Println("GetAuth, profile expired ", err)
			}

			var pf mysqlconn.Profile
			err = json.Unmarshal([]byte(p), &pf)
			if err != nil {
				fmt.Println("GetAuth, parse json fail ", err)
			}

			authres.Nickname = pf.Nickname
			authres.Picture = pf.Picture
		}
	}

	buff,_ := protocol.Marshall(&authres)

	conn.Write(buff)
	// fmt.Println("after answer ", authres, buff)
}

func onUpdateNickname(req protocol.UpdateNickname, conn net.Conn) {
	err := mysqlconn.UpdateNickname(req.Username, req.Nickname)
	if err != nil {
		fmt.Println("onUpdateNickname failed, ", err)
		return;
	}

	profilestr, err := redisconn.GetProfile(req.Username)
	if err != nil || len(profilestr) == 0 {
		fmt.Println("onUpdateNickname no profile in redis, ", err)
		return
	}

	var profile mysqlconn.Profile
	err = json.Unmarshal([]byte(profilestr), &profile)
	if err != nil {
		fmt.Println("onUpdateNickname json marshall failed, ", err)
		return;
	}

	profile.Nickname = req.Nickname

	storeProfile(&profile)

	response := protocol.UpdateNicknameRes{ 
		Username: req.Username, 
		Nickname: req.Nickname, 
		Picture: profile.Picture, 
		Reqid: req.Reqid}

	res, err := protocol.Marshall(&response)
	if err != nil {
		fmt.Println("onUpdateNickname marshall response fail, ", err)
		return
	}
	conn.Write(res)
}


func onUpdatePicture(req protocol.UpdatePicture, conn net.Conn) {
	err := mysqlconn.UpdatePicture(req.Username, req.Picture)
	if err != nil {
		fmt.Println("onUpdatePicture failed, ", err)
		return;
	}

	profilestr, err := redisconn.GetProfile(req.Username)
	if err != nil || len(profilestr) == 0 {
		fmt.Println("onUpdatePicture no profile in redis, ", err)
		return
	}

	var profile mysqlconn.Profile
	err = json.Unmarshal([]byte(profilestr), &profile)
	if err != nil {
		fmt.Println("onUpdatePicture json marshall failed, ", err)
		return;
	}

	profile.Picture = req.Picture

	storeProfile(&profile)

	response := protocol.UpdatePictureRes{ 
		Username: req.Username, 
		Nickname: profile.Nickname,
		Picture: req.Picture, 
		Reqid: req.Reqid}

	res, err := protocol.Marshall(&response)
	if err != nil {
		fmt.Println("onUpdatePicture marshall response fail, ", err)
		return
	}
	conn.Write(res)
}

func onGetProfile(req protocol.GetProfile, conn net.Conn) {
	var pf mysqlconn.Profile
	profilestr, err := redisconn.GetProfile(req.Username)
	if err != nil || len(profilestr) == 0 {
		fmt.Println("onUpdatePicture no profile in redis, ", err)

		
		pf, err = mysqlconn.QueryProfile(req.Username)
		if err != nil {
			fmt.Println("onGetProfile ", err)
			return
		}
		storeProfile(&pf)
	} else {
		err = json.Unmarshal([]byte(profilestr), &pf)
		if err != nil {
			fmt.Println("onGetProfile ", err)
			return
		}
	}

	response := protocol.GetProfileRes{ 
		Username: req.Username, 
		Nickname: pf.Nickname,
		Picture: pf.Picture,
	}

	res, err := protocol.Marshall(&response)
	if err != nil {
		fmt.Println("onUpdatePicture marshall response fail, ", err)
		return
	}
	conn.Write(res)
}


func storeProfile(p *mysqlconn.Profile) {
	redisconn.SetPassword(p.Username, p.Password)

	buff, err := json.Marshal(p)
	if err != nil {
		fmt.Println("storeProfile to redis fail, ", err)
		return
	}
	redisconn.SetProfile(p.Username, string(buff))
}

func Handle_conn(conn net.Conn) {
	fmt.Println("enter handler conn")
	uri, sockbuff := doRead(conn)

	fmt.Println("uri: ", uri)

	switch uri {
	case 1:
		var request  protocol.AuthInfo
		err := protocol.Unmarshall(&request, sockbuff)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(request.Username, " ", request.Password, " ", request.Reqid)
		onAuthInfo(request, conn)
	
	case 3: 
	    var request  protocol.UpdateNickname
		err := protocol.Unmarshall(&request, sockbuff)
		if err != nil {
			fmt.Println(err)
		}

		onUpdateNickname(request, conn)
	case 5:
		var request  protocol.UpdatePicture
		err := protocol.Unmarshall(&request, sockbuff)
		if err != nil {
			fmt.Println(err)
		}

		onUpdatePicture(request, conn)
	case 7:
		var request protocol.GetProfile
		err := protocol.Unmarshall(&request, sockbuff)
		if err != nil {
			fmt.Println(err)
		}
		onGetProfile(request, conn)
	default:
		fmt.Println("unsupport uri")

	}
}

func main() {
    addr := "0.0.0.0:3557" //表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
    listener,err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

	redisconn.InitPool()
    for {
        conn,err := listener.Accept() //用conn接收链接
        if err != nil {
            log.Fatal(err)
    	}
        go Handle_conn(conn)  //开启多个协程。
	}
}