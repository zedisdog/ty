package pay

import (
	"fmt"
	"strings"

	"github.com/zedisdog/ty/errx"
)

func BuildAttachmentsMap(attachments map[string]string) (attach string) {
	if len(attachments) < 1 {
		return
	}
	result := make([]string, 0, 10)
	for key, value := range attachments {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}

	return BuildAttachmentsSlice(result)
}

func BuildAttachmentsSlice(attachments []string) string {
	return strings.Join(attachments, "&")
}

func ParseAttachmentsSlice(attachments string) []string {
	return strings.Split(attachments, "&")
}

func ParseAttachmentsMap(attachments string) (result map[string]string, err error) {
	result = make(map[string]string, 10)
	s := ParseAttachmentsSlice(attachments)
	for _, item := range s {
		if !strings.Contains(item, "=") {
			return nil, errx.New(fmt.Sprintf("there is no key value pair exists in <%s>", item))
		}
		ss := strings.Split(item, "=")
		if len(ss) != 2 {
			return nil, errx.New(fmt.Sprintf("there is multi '=' in <%s>", item))
		}
		result[ss[0]] = ss[1]
	}
	return
}
