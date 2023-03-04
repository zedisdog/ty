package code

type code struct {
	code    int
	message string
	detail  map[string]interface{}
}

func (c code) Code() int {
	return c.code
}

func (c code) Message() string {
	return c.message
}

func (c code) Detail() map[string]interface{} {
	return c.detail
}
