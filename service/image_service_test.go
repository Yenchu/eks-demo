package service_test

import (
	"eks-demo/service"
	"io/ioutil"
	"os"
	"testing"
)

func TestResizeImage(t *testing.T) {

	bucket := os.Getenv("S3_BUCKET")
	key := "test.jpg"

	imageSvc := service.NewImageService()

	err := imageSvc.ResizeImage(bucket, key)
	if err != nil {
		t.Fatalf("Resize failed: %v", err)
	}
}

func TestResize(t *testing.T) {

	width := 2048
	height := 1024
	imageType := "jpg"

	file, err := os.Open("/Users/yenchu/data/serverless-demo/IMG_20180904_101908.jpg")
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}

	imageSvc := service.NewImageService()

	resizedImg, err := imageSvc.Resize(file, imageType, width, height)
	if err != nil {
		t.Fatalf("Resize failed: %v", err)
	}

	err = ioutil.WriteFile("/Users/yenchu/data/serverless-demo/resized.jpg", resizedImg, 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}
