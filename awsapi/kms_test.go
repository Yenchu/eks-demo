package awsapi_test

import (
	"eks-demo/awsapi"
	"os"
	"testing"
)

func newTestKmsAPI() *awsapi.KmsAPI {

	keyID := os.Getenv("KMS_KEY_ID")

	cfg := &awsapi.KmsConfig{
		KeyID: keyID,
	}

	return awsapi.NewKmsAPI(cfg)
}

func TestListKeys(t *testing.T) {

	kmsApi := newTestKmsAPI()

	resp, err := kmsApi.ListKeys()
	if err != nil {
		t.Fatalf("ListKeys failed: %v", err)
	}
	t.Logf("ListKeys: %v", resp)
}

func TestEncrypt(t *testing.T) {

	content := []byte("plain data")

	kmsApi := newTestKmsAPI()

	resp, err := kmsApi.Encrypt(content)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	t.Logf("Encrypt: %v", resp)
}
