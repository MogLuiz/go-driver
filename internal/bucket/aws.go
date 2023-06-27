package bucket

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsConfig struct {
	Config         *aws.Config
	BucketDownload string
	BucketUpload   string
}

func newAwsSession(cfg AwsConfig) *awsSession {
	return &awsSession{
		sess:           session.New(cfg.Config),
		bucketDownload: cfg.BucketDownload,
		bucketUpload:   cfg.BucketUpload,
	}
}

type awsSession struct {
	sess           *session.Session
	bucketDownload string
	bucketUpload   string
}

func (as *awsSession) Upload(r io.Reader, key string) error {
	return nil
}

func (as *awsSession) Download(src string, dst string) (file *os.File, err error) {
	file, err = os.Create(dst)
	if err != nil {
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(as.sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(as.bucketDownload),
			Key:    aws.String(src),
		})

	return
}

func (as *awsSession) Delete(key string) error {
	return nil
}