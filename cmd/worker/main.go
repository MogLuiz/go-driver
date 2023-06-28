package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/MogLuiz/go-driver/internal/bucket"
	"github.com/MogLuiz/go-driver/internal/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func main() {
	rabbitmqConfig := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBIT_URL"),
		TopicName: os.Getenv("RABBIT_TOPIC_NAME"),
		TimeOut:   time.Second * 30,
	}

	// create new Queue
	queueConnection, err := queue.New(queue.RabbitMQ, rabbitmqConfig)
	if err != nil {
		panic(err)
	}

	// create channel to consume messages
	c := make(chan queue.QueueDTO)
	queueConnection.Consume(c)

	// bucket config
	bucketConfig := bucket.AwsConfig{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), ""),
		},
		BucketDownload: "golang-drive-raw",
		BucketUpload:   "golang-drive-gzip",
	}

	// create new Bucket session
	bucket, err := bucket.New(bucket.AwsProvider, bucketConfig)
	if err != nil {
		panic(err)
	}

	for msg := range c {
		src := fmt.Sprintf("%s/%s", msg.Path, msg.Filename)
		dst := fmt.Sprintf("%d_%s", msg.ID, msg.Filename)

		file, err := bucket.Download(src, dst)
		if err != nil {
			log.Printf("Error to download file %s: %s", src, err.Error())
			continue
		}

		// read file
		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error to read file %s: %s", src, err.Error())
			continue
		}

		// compress file
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err = zw.Write(body)
		if err != nil {
			log.Printf("Error to compress file %s: %s", src, err.Error())
			continue
		}

		if err := zw.Close(); err != nil {
			log.Printf("Error to close compress file %s: %s", src, err.Error())
			continue
		}

		// upload file
		zr, err := gzip.NewReader(&buf)
		if err != nil {
			log.Printf("Error to create gzip reader: %s", err.Error())
			continue
		}

		err = bucket.Upload(zr, src)
		if err != nil {
			log.Printf("Error to upload file %s: %s", src, err.Error())
			continue
		}

		// delete file
		err = os.Remove(dst)
		if err != nil {
			log.Printf("Error to delete file %s: %s", dst, err.Error())
			continue
		}
	}
}
