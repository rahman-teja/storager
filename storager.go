package storager

import (
	"context"
	"io"
)

type Storage interface {
	GetObject(ctx context.Context, token, filename string, opts GetOptions) (reader io.Reader, err error)
}

type WriteableStorage interface {
	Storage
	PutObject(ctx context.Context, token, filename string, reader io.ReadSeeker, objectSize int64, opts PutOptions) (id string, err error)
	DeleteObject(ctx context.Context, token, filename string, opts RemoveOptions) (err error)
}
