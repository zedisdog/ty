package strings

import (
	"fmt"
	"net/url"
	"strings"
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
	tmp, err := url.Parse(u)
	if err == nil {
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
