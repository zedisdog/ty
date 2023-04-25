package strings

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func ContainersAny(str string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(str, substr) {
			return true
		}
	}

	return false
}

func EncodeQuery(u string) string {
	if strings.Contains(u, "://") {
		tmp, err := url.Parse(u)
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%s://%s@%s%s?%s", tmp.Scheme, tmp.User, tmp.Host, tmp.Path, tmp.Query().Encode())
	} else {
		urlArr := strings.Split(u, "?")
		if len(urlArr) > 1 {
			q, err := url.ParseQuery(urlArr[1])
			if err != nil {
				panic(err)
			}
			return fmt.Sprintf("%s?%s", urlArr[0], q.Encode())
		} else {
			return u
		}
	}
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func RandNumeric(len int) string {
	numeric := [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	r := rand.New(rand.NewSource(time.Now().Unix()))

	var sb strings.Builder

	for i := 0; i < len; i++ {
		sb.WriteByte(numeric[r.Intn(10)])
	}

	return sb.String()
}
