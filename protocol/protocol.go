package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type ProtocolBase interface {
	Uri() int32
}


type AuthInfo struct {
	Username string
	Token string
	Reqid    int
}

func (this AuthInfo) Uri() int32 {
	return 1;
}

type AuthResult struct {
	Username string
	Nickname string
	Picture string
	Result   bool
	Reqid    int
}

func (this AuthResult) Uri() int32 {
	return 2;
}

type UpdateNickname struct {
	Username string
	Nickname string
	Reqid    int32
}

func (this UpdateNickname) Uri() int32 {
	return 3
}

type UpdateNicknameRes struct {
	Username string
	Nickname string
	Picture string
	Reqid    int32
}

func (this UpdateNicknameRes) Uri() int32 {
	return 4
}

type UpdatePicture struct {
	Username string
	Picture string
	Reqid    int32
}

func (this UpdatePicture) Uri() int32 {
	return 5
}

type UpdatePictureRes struct {
	Username string
	Nickname string
	Picture string
	Reqid    int32
}

func (this UpdatePictureRes) Uri() int32 {
	return 6
}

type GetProfile struct {
	Username string
}

func (this GetProfile) Uri() int32 {
	return 7
}

type GetProfileRes struct {
	Username string
	Nickname string
	Picture string
}

func (this GetProfileRes) Uri() int32 {
	return 8
}

func Marshall(data ProtocolBase) ([]byte, error) {
	buf := new(bytes.Buffer)

	b, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	binary.Write(buf, binary.LittleEndian, data.Uri())
	var len int32 = int32(len(b));
	binary.Write(buf, binary.LittleEndian, len);
	binary.Write(buf, binary.LittleEndian, b);

	fmt.Println("Marshall uri ", data.Uri(), " len ", len)

	return buf.Bytes(), nil	
}

func Unmarshall(data ProtocolBase, buff []byte) error {
	return json.Unmarshal(buff, &data)
}
