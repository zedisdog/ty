package common

const (
	ErrBusy        = -1
	ErrSuccess     = 0
	ErrInvalidCode = 40029
	ErrHighRisk    = 40226
)

type ErrorResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
