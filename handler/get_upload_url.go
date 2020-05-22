package handler

import (
	"eks-demo/model"
	"eks-demo/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var bucket string
var uploadSvc *service.UploadService

func InitUpload() {

	bucket = os.Getenv("S3_BUCKET")
	log.WithFields(log.Fields{"bucketName": bucket}).Info("S3_BUCKET")

	uploadSvc = service.NewUploadService()
}

func GetUploadURL(c *gin.Context) {

	reqData := &model.GetUploadURLRequest{}

	if err := c.ShouldBindJSON(reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.WithFields(log.Fields{"file": reqData.File, "contentType": reqData.ContentType,
		"width": reqData.Width, "height": reqData.Height}).Info("Get upload url request")

	reqData.Bucket = bucket

	respData, err := uploadSvc.GetUploadURL(reqData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, respData)
}
