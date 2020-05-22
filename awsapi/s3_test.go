package awsapi_test

import (
	"os"
	"serverless-demo/awsapi"
	"testing"
	"time"
)

func TestGetPutObjectPreSignURL(t *testing.T) {

	metadata := map[string]string{
		"width": "800",
		"height": "600",
	}

	reqData := &awsapi.S3PreSignURLRequest{
		Bucket: os.Getenv("S3_BUCKET"),
		Key: "resize/test2.jpg",
		ContentType: "image/jpg",
		Metadata: metadata,
		TTL: 10 * time.Minute,
	}

	s3Api := awsapi.NewS3API()

	signedURL, err := s3Api.GetPutObjectPreSignURL(reqData)
	if err != nil {
		t.Fatalf("GetPutObjectPreSignURL failed: %v", err)
	}
	t.Logf("GetPutObjectPreSignURL:\n%s", signedURL)
}

func TestGetPutObjectPreSignURLHeaders(t *testing.T) {

	metadata := map[string]string{
		"width": "800",
		"height": "600",
	}

	reqData := &awsapi.S3PreSignURLRequest{
		Bucket: os.Getenv("S3_BUCKET"),
		Key: "resize/test.jpg",
		ContentType: "image/jpg",
		Metadata: metadata,
		TTL: 10 * time.Minute,
	}

	s3Api := awsapi.NewS3API()

	url, headers, err := s3Api.GetPutObjectPreSignURLHeaders(reqData)
	if err != nil {
		t.Fatalf("GetPutObjectPreSignURLHeaders failed: %v", err)
	}
	t.Logf("GetPutObjectPreSignURLHeaders:\nurl=%s\nheaders=%v", url, headers)
}

func TestPutObject(t *testing.T) {

	image, err := os.Open("/Users/yenchu/data/serverless-demo/IMG_20180904_101908.jpg")
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}

	bucket := os.Getenv("S3_BUCKET")
	key := "resize/test.jpg"
	contentType := "image/jpg"

	metadata := map[string]string{
		"width": "1024",
		"height": "1024",
	}

	s3Api := awsapi.NewS3API()

	resp, err := s3Api.PutObject(image, bucket, key, contentType, metadata)
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}
	t.Logf("PutObject result: %+v", resp)
}