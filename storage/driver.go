package storage

type IDriver interface {
	Put(path string, data []byte) error
	Get(path string) ([]byte, error)
	Remove(path string) error
}

type IDriverHasMime interface {
	MimeType(path string) string
}

type IDriverHasPath interface {
	//Path 绝对路径
	Path(path string) string
}

type IDriverHasBase64 interface {
	Base64(path string) (string, error)
}

type IDriverCanGetSize interface {
	Size(path string) (int, error)
}

type IDriverHasUrl interface {
	Url(path string) string
}

type IDriverCanAppend interface {
	Append(path string, data []byte) (err error)
}
