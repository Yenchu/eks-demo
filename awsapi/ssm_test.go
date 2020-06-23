package awsapi_test

import (
	"eks-demo/awsapi"
	"testing"
)

func TestGetDecryptedParameter(t *testing.T) {

	ssmApi := awsapi.NewSsmAPI()

	val, err := ssmApi.GetDecryptedParameter("/applications/ServerlessDemo/CloudFront/PrivateKey")
	if err != nil {
		t.Fatalf("GetDecryptedParameter failed: %v", err)
	}
	t.Logf("GetDecryptedParameter: %v", val)

}

func TestGetParameters(t *testing.T) {

	ssmApi := awsapi.NewSsmAPI()

	vals, err := ssmApi.GetParameters("/applications/ServerlessDemo/CloudFront/KeyId", "/applications/ServerlessDemo/CloudFront/PrivateKey")
	if err != nil {
		t.Fatalf("GetParameters failed: %v", err)
	}
	t.Logf("GetParameters: %v", vals)

}
