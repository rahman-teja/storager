package storager

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type minioStorage struct {
	client     MinioClient
	region     string
	bucketName string
}

func NewMinioStorage(client *minio.Client, region, bucketName string) WriteableStorage {
	ret := &minioStorage{
		client:     client,
		region:     region,
		bucketName: bucketName,
	}

	err := ret.mustCreateBucket(context.Background())
	if err != nil {
		return nil
	}

	return ret
}

func (m minioStorage) GetObject(ctx context.Context, token, filename string, opts GetOptions) (reader io.Reader, err error) {
	var respGetObj *minio.Object

	respGetObj, err = m.client.GetObject(ctx, m.bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	defer respGetObj.Close()

	_, err = respGetObj.Stat()
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(respGetObj)

	reader = buf
	return
}

func (m minioStorage) PutObject(ctx context.Context, token, filename string, reader io.ReadSeeker, objectSize int64, opts PutOptions) (id string, err error) {
	var uplInfo minio.UploadInfo

	putOpts := minio.PutObjectOptions{
		ContentType:             opts.ContentType,
		UserMetadata:            opts.UserMetadata,
		UserTags:                opts.UserTags,
		ContentEncoding:         opts.ContentEncoding,
		ContentDisposition:      opts.ContentDisposition,
		ContentLanguage:         opts.ContentLanguage,
		CacheControl:            opts.CacheControl,
		WebsiteRedirectLocation: opts.WebsiteRedirectLocation,
	}

	uplInfo, err = m.client.PutObject(ctx, m.bucketName, filename, reader, objectSize, putOpts)
	if err != nil {
		return
	}

	id = uplInfo.Key
	return
}

func (m minioStorage) DeleteObject(ctx context.Context, token, filename string, opts RemoveOptions) (err error) {
	rmvOpts := minio.RemoveObjectOptions{
		VersionID: opts.VersionID,
	}

	err = m.client.RemoveObject(ctx, m.bucketName, filename, rmvOpts)
	return
}

func (m minioStorage) mustCreateBucket(ctx context.Context) (err error) {
	var exists bool

	exists, err = m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return
	}

	if exists {
		return
	}

	opts := minio.MakeBucketOptions{
		Region: m.region,
	}

	err = m.client.MakeBucket(ctx, m.bucketName, opts)
	return
}
