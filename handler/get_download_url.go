package handler

import (
	"eks-demo/model"
	"eks-demo/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const DefaultHTTPScheme = "https"

var cfDomain string
var downloadSvc *service.DownloadService

func InitDownload() error {

	cfDomain = os.Getenv("CF_DOMAIN_NAME")
	log.WithFields(log.Fields{"domainName": cfDomain}).Info("CF_DOMAIN_NAME")

	svc, err := service.NewDownloadService()
	if err != nil {
		return err
	}

	downloadSvc = svc
	return nil
}

func GetDownloadURL(c *gin.Context) {

	reqData := &model.GetDownloadURLRequest{}

	if err := c.ShouldBindJSON(reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.WithFields(log.Fields{"file": reqData.File}).Info("Get download url request")

	reqData.Scheme = DefaultHTTPScheme
	reqData.Domain = cfDomain

	respData, err := downloadSvc.GetDownloadURL(reqData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, respData)
}
