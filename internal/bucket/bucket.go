package bucket

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

type BucketType int

const (
	AwsProvider BucketType = iota
)

func New(bt BucketType, cfg any) (b *Bucket, err error) {
	rt := reflect.TypeOf(cfg)

	switch bt {
	case AwsProvider:
		// TODO implement AWS provider
	default:
		return nil, fmt.Errorf("Bucket type not implemented")
	}

	return
}

type BucketInterface interface {
	Upload(io.Reader, string) error
	Download(string, string) (*os.File, error)
	Delete(string) error
}

type Bucket struct {
	p BucketInterface
}

func (b *Bucket) Upload(r io.Reader, key string) error {
	return b.p.Upload(r, key)
}

func (b *Bucket) Download(src string, dst string) (file *os.File, err error) {
	return b.p.Download(src, dst)
}

func (b *Bucket) Delete(key string) error {
	return b.p.Delete(key)
}
