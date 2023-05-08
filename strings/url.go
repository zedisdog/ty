package strings

import (
	"fmt"
	"net/url"
	"strings"
)

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
