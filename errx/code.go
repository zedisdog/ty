package errx

import "github.com/zedisdog/ty/i18n"

const (
	Nil Code = -1
)

var codeMap = map[Code]string{
	Nil: "nil code",
}

type Code int

// IsPreDefined reports if the code is predefined
func (code Code) IsPreDefined() bool {
	for c := range codeMap {
		if c == code {
			return true
		}
	}

	return false
}

// Message gets meaning of code
//
//	Returns: empty means the code is none predefined code
func (code Code) Message(lang ...string) string {
	l := i18n.DefaultLang
	if len(lang) > 0 {
		l = lang[0]
	}

	text := codeMap[code]

	msg := i18n.TranslateByLang(l, text)

	if msg == "" {
		return codeMap[code]
	} else {
		return msg
	}
}

// Register registers the code as predefined code
func Register(code Code, meaning string) (err error) {
	if code.IsPreDefined() {
		err = New("code already exists")
		return
	}

	if meaning == "" {
		err = New("meaning is required")
		return
	}

	codeMap[code] = meaning

	return
}
