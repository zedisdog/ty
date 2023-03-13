package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/zedisdog/ty/errx"
)

//go:generate mockgen -destination=./test/storage.go -package=test github.com/zedisdog/sweetbean/storage IStorage

var _ IStorage = (*Storage)(nil)

type IStorage interface {
	IDriver
	PutString(path string, data string) error
	GetString(path string) (data string, err error)
	PutFile(path string, file *multipart.FileHeader) (err error)
	MimeType(path string) (string, error)
	Path(path string) (string, error)
	Base64(path string) (string, error)
	Size(path string) (int, error)
	Url(path string) (string, error)
	PutFileQuick(file *multipart.FileHeader, directory string) (path string, err error)
	PutFileBytesQuick(data []byte, ext string, dir string) (path string, err error)
	Append(path string, data []byte) (err error)
}

func NewStorage(driver IDriver) *Storage {
	return &Storage{
		driver: driver,
	}
}

type Storage struct {
	driver IDriver
}

func (s Storage) randFileName(ext string) (name string, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	if strings.HasPrefix(ext, ".") {
		return fmt.Sprintf("%s%s", id.String(), ext), nil
	} else {
		return fmt.Sprintf("%s.%s", id.String(), ext), nil
	}
}

// PutFileQuick is similar than PutFile, but don't set filename.
func (s Storage) PutFileQuick(file *multipart.FileHeader, directory string) (path string, err error) {
	fileName, err := s.randFileName(filepath.Ext(file.Filename))
	if err != nil {
		return
	}
	//path = directory/xxxx.jpg
	path = fmt.Sprintf(
		"%s%s%s",
		strings.Trim(directory, "\\/"),
		string(os.PathSeparator),
		fileName,
	)
	err = s.PutFile(path, file)
	return
}

func (s Storage) PutFileBytesQuick(data []byte, ext string, dir string) (path string, err error) {
	fileName, err := s.randFileName(ext)
	if err != nil {
		return
	}
	path = fmt.Sprintf(
		"%s%s%s",
		strings.Trim(dir, "\\/"),
		string(os.PathSeparator),
		fileName,
	)
	err = s.Put(path, data)
	return
}

func (s Storage) Put(path string, data []byte) error {
	return s.driver.Put(path, data)
}

func (s Storage) Get(path string) ([]byte, error) {
	return s.driver.Get(path)
}

func (s Storage) Remove(path string) error {
	return s.driver.Remove(path)
}

func (s Storage) PutString(path string, data string) error {
	return s.driver.Put(path, []byte(data))
}

func (s Storage) GetString(path string) (data string, err error) {
	tmp, err := s.driver.Get(path)
	if err != nil {
		return
	}
	data = string(tmp)
	return
}

func (s Storage) PutFile(path string, file *multipart.FileHeader) (err error) {
	fp, err := file.Open()
	if err != nil {
		return
	}
	defer func() {
		_ = fp.Close()
	}()
	data, err := io.ReadAll(fp)
	if err != nil {
		return
	}
	return s.Put(path, data)
}

func (s Storage) MimeType(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(IDriverHasMime); ok {
		return ss.MimeType(path), nil
	}
	return "", errx.New("driver is not implement interface <IDriverHasMime>")
}

func (s Storage) Path(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(IDriverHasPath); ok {
		return ss.Path(path), nil
	}
	return "", errx.New("driver is not implement interface <IDriverHasPath>")
}

func (s Storage) Base64(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(IDriverHasBase64); ok {
		return ss.Base64(path)
	}
	return "", errx.New("driver is not implement interface <IDriverHasBase64>")
}

func (s Storage) Size(path string) (int, error) {
	if ss, ok := interface{}(s.driver).(IDriverCanGetSize); ok {
		return ss.Size(path)
	}
	return 0, errx.New("driver is not implement interface <IDriverCanGetSize>")
}

func (s Storage) Url(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(IDriverHasUrl); ok {
		return ss.Url(path), nil
	}
	return "", errx.New("driver is not implement interface <IDriverHasUrl>")
}

func (s Storage) Append(path string, data []byte) (err error) {
	if ss, ok := interface{}(s.driver).(IDriverCanAppend); ok {
		return ss.Append(path, data)
	}
	return errx.New("driver is not implement interface <IDriverCanAppend>")
}
