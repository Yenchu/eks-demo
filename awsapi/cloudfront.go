package awsapi

import (
	"crypto/rsa"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/sign"
	"time"
)

func NewCloudFrontAPI(cfg *CloudFrontConfig) *CloudFrontAPI {

	client := cloudfront.New(LoadAWSConfig())

	return &CloudFrontAPI{
		cfg:    cfg,
		client: client,
	}
}

type CFSignURLRequest struct {
	URL string
	TTL time.Duration
}

type CloudFrontConfig struct {
	KeyID      string
	PrivateKey *rsa.PrivateKey
}

type CloudFrontAPI struct {
	cfg    *CloudFrontConfig
	client *cloudfront.Client
}

func (api *CloudFrontAPI) GetSignURL(reqData *CFSignURLRequest) (string, error) {

	signer := sign.NewURLSigner(api.cfg.KeyID, api.cfg.PrivateKey)

	return signer.Sign(reqData.URL, time.Now().UTC().Add(reqData.TTL))
}
