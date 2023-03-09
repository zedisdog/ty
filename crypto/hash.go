package crypto

import (
	"crypto"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// Sha1 sha1算法
func Sha1(str string) string {
	hash := crypto.SHA1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func CheckSha1(expect string, actual []byte) bool {
	encoder := sha1.New()
	encoder.Write(actual)
	test := fmt.Sprintf("%x", encoder.Sum(nil))
	return expect == test
}
