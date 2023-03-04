package code

type ICode interface {
	Code() int
	Message() string
	Detail() map[string]interface{}
}

var (
	Nil      = &code{-1, "no code", nil}
	NotFound = &code{404, "not found", nil}
)
