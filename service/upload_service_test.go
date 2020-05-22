package service_test

import (
	"eks-demo/model"
	"eks-demo/service"
	"testing"
)

func TestGetUploadURL(t *testing.T) {

	req := &model.GetUploadURLRequest{
		Bucket:      "abc",
		File:        "test.jpg",
		ContentType: "image/jpg",
		Width:       800,
		Height:      600,
	}

	uploadSvc := service.NewUploadService()

	resp, err := uploadSvc.GetUploadURL(req)
	if err != nil {
		t.Fatalf("GetUploadURL failed: %v", err)
	}
	t.Logf("GetUploadURL: %v", resp)
}
