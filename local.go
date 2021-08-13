package storager

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
)

type localStorage struct {
	basepath string
}

func NewLocalStorage(basepath string) WriteableStorage {
	return &localStorage{
		basepath: basepath,
	}
}

func (l localStorage) GetObject(ctx context.Context, token, filename string, opts GetOptions) (reader io.Reader, err error) {
	var file *os.File
	var path string = filepath.Join(l.basepath, filename)

	file, err = os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(file)

	reader = buf
	return
}

func (l localStorage) PutObject(ctx context.Context, token, filename string, reader io.ReadSeeker, objectSize int64, opts PutOptions) (s string, err error) {
	var file *os.File
	var path string = filepath.Join(l.basepath, filename)

	file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return
	}

	s = path
	return
}

func (l localStorage) DeleteObject(ctx context.Context, token, filename string, opts RemoveOptions) (err error) {
	var path string = filepath.Join(l.basepath, filename)

	err = os.Remove(path)
	return
}
