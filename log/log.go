package log

func NewField(name string, value interface{}) *Field {
	return &Field{
		Name:  name,
		Value: value,
	}
}

type Field struct {
	Name  string
	Value interface{}
}

type Level int

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

type ILog interface {
	Trace(msg string, fields ...*Field)
	Debug(msg string, fields ...*Field)
	Info(msg string, fields ...*Field)
	Warn(msg string, fields ...*Field)
	Error(msg string, fields ...*Field)
	Fatal(msg string, fields ...*Field)
	Log(msg string, level Level, fields ...*Field)
}
