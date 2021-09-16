package rpcclient

import (
	"bytes"

	"encoding/binary"
	"encoding/json"
	"fmt"

	"entry_task/protocol"
	"net"
)

var reqid int = 0;

func getReqid() int {
	defer increaseReqid()
	return reqid
}

func increaseReqid() {
	reqid++
}

func GetAuth(username string, password string) (bool, string, string, string, error) {
	// connectServer
	conn, err := net.Dial("tcp", ":3557")
    if err != nil {
        fmt.Println(err)
		return false, "", "", "", err
    }
	defer conn.Close()

	// new Message
	authinfo := protocol.AuthInfo{ 
		Username: username, 
		Password: password, 
		Reqid: reqid,
	}
	//Marshal
	b, err := protocol.Marshall(&authinfo)
	if err != nil {
		fmt.Println(err)
		return false, "", "", "", err
	}
	conn.Write(b)

	// wait for answer
	uri, res := doRead(conn)
	var authres protocol.AuthResult
	if uri == 2 {
		err = protocol.Unmarshall(&authres, res)
		if err != nil {
			return false, "", "", "", err
		}
	}
    fmt.Println("GetAuth, ", string(res))

	var response protocol.AuthResult
	err = json.Unmarshal(res, &response)
	
	if err != nil {
		fmt.Println("GetAuth ", err)
		return false, "", "", "", err
	}

	return response.Result, response.Username, response.Nickname, response.Picture, nil
}

func UpdateNickname(username string, nickname string) (string, string, string, error) {
	conn, err := net.Dial("tcp", ":3557")
    if err != nil {
        fmt.Println(err)
		return "", "", "", err
    }
	defer conn.Close()

	msg := protocol.UpdateNickname{ Username: username, Nickname: nickname}
	b, err := protocol.Marshall(&msg)

	if err != nil {
		fmt.Println(err)
		return "", "", "", err
	}
	conn.Write(b)

	_, res := doRead(conn)

	// onAnswer()
	var response protocol.UpdateNicknameRes

	err = protocol.Unmarshall(&response, res)
	if err != nil {
		return "", "", "", err
	}

    fmt.Println("UpdateNickname, ", string(res))
	return response.Username, response.Nickname, response.Picture, nil
}

func UpdatePicture(username string, Picture string)  (string, string, string, error) {
	conn, err := net.Dial("tcp", ":3557")
    if err != nil {
        fmt.Println(err)
		return "", "", "", err
    }
	defer conn.Close()

	msg := protocol.UpdatePicture{
		Username: username,
		Picture: Picture,
	}
	b, err := protocol.Marshall(&msg)

	if err != nil {
		fmt.Println(err)
		return "", "", "", err
	}
	conn.Write(b)


	uri, res := doRead(conn)

	// onAnswer()
	var response protocol.UpdatePictureRes
	if uri == 6 {
		err = protocol.Unmarshall(&response, res)
		if err != nil {
			return "", "", "", err
		}
	}
    fmt.Println("updatePicture, ", string(res))
	return response.Username, response.Nickname, response.Picture, nil
}

func GetProfile(username string) (string, string, string, error) {
	conn, err := net.Dial("tcp", ":3557")
    if err != nil {
        fmt.Println(err)
		return "", "", "", err
    }
	defer conn.Close()

	msg := protocol.GetProfile{
		Username: username,
	}
	b, err := protocol.Marshall(&msg)

	if err != nil {
		fmt.Println(err)
		return "", "", "", err
	}
	conn.Write(b)


	_, res := doRead(conn)

	// onAnswer()
	var response protocol.GetProfileRes
	err = protocol.Unmarshall(&response, res)
	if err != nil {
		return "", "", "", err
	}
    fmt.Println("GetProfile, ", string(res))
	return response.Username, response.Nickname, response.Picture, nil	
}

func peekhead(buff []byte) (uri int32, length int32, e error) {
	decbuffer := bytes.NewBuffer(buff)

	fmt.Println("peek head: ", buff);

	fmt.Println("peek head: ", decbuffer);
	err := binary.Read(decbuffer, binary.LittleEndian, &uri)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(decbuffer, binary.LittleEndian, &length)
	if err != nil {
		fmt.Println(err)
	}

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
		return 0, nil
	}

	fmt.Println("uri: ", uri, " length: ", length, "  n: ", n, " err: ", err)
	fmt.Println(buf)

	if n < int(length + 8) {
		return 0, nil
	}

	return uri, buf[8:8+length]
}

// func main() {
// 	GetAuth("wanlb", "123456");
// 	UpdateNickname("wanlb", "Rubin123")
// 	UpdatePicture("wanlb", "Rubin123")
// }