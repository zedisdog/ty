package errx

const (
	Nil      Code = -1
	NotFound Code = 404
)

var codeMap = map[Code]string{
	Nil:      "nil code",
	NotFound: "not found",
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
func (code Code) Message() string {
	return codeMap[code]
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
