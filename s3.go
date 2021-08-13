package storager

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Storage struct {
	conf       *aws.Config
	bucketName string
}

func NewS3Storage(conf *aws.Config, bucketName string) WriteableStorage {
	return &s3Storage{
		conf:       conf,
		bucketName: bucketName,
	}
}

func (l s3Storage) createNewService() (svc *s3.S3) {
	sess := session.Must(session.NewSession(
		l.conf,
	))

	svc = s3.New(sess)
	return
}

func (l s3Storage) GetObject(ctx context.Context, token, filename string, opts GetOptions) (reader io.Reader, err error) {
	svc := l.createNewService()

	resp, err := svc.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(l.bucketName),
			Key:    aws.String(filename),
		},
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	reader = buf
	return
}

func (l s3Storage) PutObject(ctx context.Context, token, filename string, reader io.ReadSeeker, objectSize int64, opts PutOptions) (s string, err error) {
	svc := l.createNewService()

	resp, err := svc.PutObjectWithContext(
		ctx,
		&s3.PutObjectInput{
			Body:                    reader,
			Bucket:                  aws.String(l.bucketName),
			ContentLength:           aws.Int64(objectSize),
			Key:                     aws.String(filename),
			ContentType:             aws.String(opts.ContentType),
			ContentDisposition:      aws.String(opts.ContentDisposition),
			CacheControl:            aws.String(opts.CacheControl),
			ContentEncoding:         aws.String(opts.ContentEncoding),
			ContentLanguage:         aws.String(opts.ContentLanguage),
			Metadata:                aws.StringMap(opts.UserMetadata),
			WebsiteRedirectLocation: aws.String(opts.WebsiteRedirectLocation),
			Expires:                 aws.Time(opts.Expires),
		},
	)

	s = awsutil.StringValue(resp)
	return
}

func (l s3Storage) DeleteObject(ctx context.Context, token, filename string, opts RemoveOptions) (err error) {
	svc := l.createNewService()

	_, err = svc.DeleteObjectWithContext(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(l.bucketName),
			Key:    aws.String(filename),
		},
	)

	return
}
