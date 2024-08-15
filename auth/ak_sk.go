package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// GenAccessKey 生成随机字符串作为 AccessKey
func GenAccessKey(userId string) string {

	// akStr := base64.URLEncoding.EncodeToString(uuid()) + ":" + userId
	return base64.URLEncoding.EncodeToString([]byte(userId))
}

// GenSecretKey 生成随机字符串作为 SecretKey
func GenSecretKey() string {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(key)
}

// uuid returns a new UUID v6 as a byte slice.
// func uuid() []byte {
// 	u := make([]byte, 16)

// 	ts := uint64(time.Now().UnixNano() / 1e6)
// 	binary.BigEndian.PutUint64(u[:8], ts)

// 	_, err := rand.Read(u[8:16])
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Set the version to 6 (6 << 4 is 01100000)
// 	u[6] = (u[6] | 0x0F) | 0x60
// 	// Set the variant to 10xx (for RFC 4122)
// 	u[8] = (u[8] & 0x3F) | 0x80

// 	return u
// }

//func main() {
//	accessKey := generateAccessKey("b745faf66f0c434ea50ae5aaed098936")
//	secretKey := generateSecretKey()
//
//	fmt.Println("AccessKey:", accessKey)
//	fmt.Println("SecretKey:", secretKey)
//}
