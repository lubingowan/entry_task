package token

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

const VERSION_LEN int = 3
const APPIDSTR_LEN int = 32

func version() string {
	return "001"
}

const cert string = "A_Typical_Cert"

func Sha1(password string) string {
	mac := hmac.New(sha1.New, []byte(cert))
	mac.Write([]byte(password))

	return string(mac.Sum(nil))
}

func GenSignature(uname string, password []byte) []byte {

	buffer := new(bytes.Buffer)
	buffer.WriteString(uname)
	buffer.Write(password)
	buffer.WriteString(cert)

	// binary.Write(buffer, binary.LittleEndian, []uint32{salt, gents, effts})

	key := []byte(cert)
	mac := hmac.New(sha1.New, key)
	mac.Write(buffer.Bytes())

	return mac.Sum(nil)
}

/* 生成token入口 */
func GenToken(uname string, password []byte) string {
	// 生成时间，盐值，有效期
	var gents uint32 = uint32(time.Now().Unix())
	rand.Seed(time.Now().UnixNano())
	var salt uint32 = rand.Uint32()
	var effts uint32 = 864000

	resbuffer := new(bytes.Buffer)

	// 序列号签名
	sigbuf := GenSignature(uname, password)
	var siglen uint16
	siglen = uint16(len(sigbuf))
	binary.Write(resbuffer, binary.LittleEndian, siglen)
	binary.Write(resbuffer, binary.LittleEndian, sigbuf)

	bytesbuffer := new(bytes.Buffer)
	bytesbuffer.WriteString(uname)
	crc32uname := crc32.ChecksumIEEE(bytesbuffer.Bytes())
	binary.Write(resbuffer, binary.LittleEndian, crc32uname)

	// 序列化 username password crc32
	bytesbuffer.Reset()
	bytesbuffer.Write(password)
	crc32password := crc32.ChecksumIEEE(bytesbuffer.Bytes())
	binary.Write(resbuffer, binary.LittleEndian, crc32password)

	// // 序列化盐值、生成时间、有效时间
	binary.Write(resbuffer, binary.LittleEndian, salt)
	binary.Write(resbuffer, binary.LittleEndian, gents)
	binary.Write(resbuffer, binary.LittleEndian, effts)

	res := version() + base64.StdEncoding.EncodeToString(resbuffer.Bytes())

	return res
}

func ParseToken(s string) (uint32, uint32, uint32, string) {
	if len(s) < VERSION_LEN {
		return 0, 0, 0, ""
	}

	ver := s[0:VERSION_LEN]
	if ver != version() {
		return 0, 0, 0, ""
	}

	encstr := s[VERSION_LEN:]
	if len(encstr) < 1 {
		return 0, 0, 0, ""
	}

	decodeBytes, err := base64.StdEncoding.DecodeString(encstr)
	if err != nil {
		fmt.Println("ParseToken decode base64 ", err)
		return 0, 0, 0, ""
	}

	decbuffer := bytes.NewBuffer(decodeBytes)

	var siglen uint16
	binary.Read(decbuffer, binary.LittleEndian, &siglen)
	sigbytes := make([]byte, siglen)
	decbuffer.Read(sigbytes)
	sigstr := string(sigbytes)

	var crc32uname, crc32password, salt, gents, effts uint32

	binary.Read(decbuffer, binary.LittleEndian, &crc32uname)
	binary.Read(decbuffer, binary.LittleEndian, &crc32password)
	binary.Read(decbuffer, binary.LittleEndian, &salt)
	binary.Read(decbuffer, binary.LittleEndian, &gents)
	binary.Read(decbuffer, binary.LittleEndian, &effts)

	fmt.Println(salt, gents, effts, sigstr)

	return salt, gents, effts, sigstr
}

func CheckToken(token string, uname string, password []byte) bool {

	_, gents, effts, sigstr := ParseToken(token)

	if gents == 0 {
		fmt.Println("parseToken failed")
		return false
	}

	// check if time expired, 检查token是否过期
	if (gents + effts) < uint32(time.Now().Unix()) {
		fmt.Println("token expired")
		return false
	}

	// check if signature valid, 检查签名是否有效
	signatureNow := string(GenSignature(uname, password))
	if sigstr == signatureNow {
		fmt.Println("token valid!!!!")
		return true
	} else {
		fmt.Println("token invalid????")
		return false
	}
}

func CheckTokenByToken(left string, right string) bool {

	_, gents, effts, sigstr := ParseToken(left)

	if gents == 0 {
		fmt.Println("parseToken failed")
		return false
	}

	// check if time expired, 检查token是否过期
	if (gents + effts) < uint32(time.Now().Unix()) {
		fmt.Println("token expired")
		return false
	}

	_, _, _, sign := ParseToken(right)
	fmt.Println(len(sigstr), "  ", len(sign), " ", sigstr == sign)

	if (sigstr == sign) {
		return true
	} else {
		return false
	}
}

// func test() {
// 	var token Token
// 	token.init("mycert_string")


// }
