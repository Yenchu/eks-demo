package service

import (
	"eks-demo/awsapi"
	"eks-demo/model"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/sign"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

const (
	SsmCFKeyID      = "/applications/ServerlessDemo/CloudFront/KeyId"
	SsmCFPrivateKey = "/applications/ServerlessDemo/CloudFront/PrivateKey"
	CFSignedUrlTTL  = 7 * 24 * time.Hour
)

func NewDownloadService() (*DownloadService, error) {

	ssmApi := awsapi.NewSsmAPI()

	cfCfg, err := createCFConfig(ssmApi)
	if err != nil {
		return nil, err
	}

	cfApi := awsapi.NewCloudFrontAPI(cfCfg)

	return &DownloadService{
		cfApi: cfApi,
	}, nil
}

func createCFConfig(ssmClient *awsapi.SsmAPI) (*awsapi.CloudFrontConfig, error) {

	chErr := ""
	pkCh := make(chan string)
	go func() {
		pkStr, err := ssmClient.GetDecryptedParameter(SsmCFPrivateKey)
		if err != nil {
			log.WithFields(log.Fields{"param": SsmCFPrivateKey, "error": err}).Error("GetDecryptedParameter failed")
			pkCh <- chErr
		} else {
			pkCh <- pkStr
		}
	}()

	idCh := make(chan string)
	go func() {
		keyID, err := ssmClient.GetParameter(SsmCFKeyID)
		if err != nil {
			log.WithFields(log.Fields{"param": SsmCFKeyID, "error": err}).Error("GetParameter failed")
			idCh <- chErr
		} else {
			idCh <- keyID
		}
	}()

	pkStr := <-pkCh
	keyID := <-idCh

	if pkStr == chErr {
		return nil, errors.Errorf("Get SSM parameter %s failed", SsmCFPrivateKey)
	}
	if keyID == chErr {
		return nil, errors.Errorf("Get SSM parameter %s failed", SsmCFKeyID)
	}

	pk, err := sign.LoadPEMPrivKey(strings.NewReader(pkStr))
	if err != nil {
		return nil, err
	}

	return &awsapi.CloudFrontConfig{
		KeyID:      keyID,
		PrivateKey: pk,
	}, nil
}

type DownloadService struct {
	cfApi *awsapi.CloudFrontAPI
}

func (svc *DownloadService) GetDownloadURL(req *model.GetDownloadURLRequest) (*model.GetDownloadURLResponse, error) {

	urlObj := &url.URL{
		Scheme: req.Scheme,
		Host:   req.Domain,
		Path:   model.ImagesFileDir + "/" + req.File,
	}

	reqData := &awsapi.CFSignURLRequest{
		URL: urlObj.String(),
		TTL: CFSignedUrlTTL,
	}

	url, err := svc.cfApi.GetSignURL(reqData)
	if err != nil {
		return nil, err
	}

	return &model.GetDownloadURLResponse{
		URL: url,
	}, nil
}
