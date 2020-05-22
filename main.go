package main

import (
	"eks-demo/handler"
	"eks-demo/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	AllSvc      = "*"
	UploadSvc   = "upload"
	DownloadSvc = "download"
	ResizeSvc   = "resize"
)

func main() {

	svc := os.Getenv("SVC")
	if svc == "" {
		svc = AllSvc
	}
	log.WithFields(log.Fields{"svc": svc}).Info("enable service")

	r := gin.Default()

	r.GET("/", handler.GetHealth)

	//r.GET("/liveness", handler.GetHealth)

	//r.GET("/readiness", handler.GetHealth)

	if svc == AllSvc || svc == UploadSvc {
		handler.InitUpload()
		// cannot omit path /upload, because ALB ingress controller doesn't support path rewrite currently
		r.POST("/upload/get-signed-url", handler.GetUploadURL)

		// for health check test
		r.GET("/upload", handler.GetHealth)
	}

	if svc == AllSvc || svc == DownloadSvc {
		err := handler.InitDownload()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("InitDownload failed")
			return
		}

		r.POST("/download/get-signed-url", handler.GetDownloadURL)

		r.GET("/download", handler.GetHealth)
	}

	if svc == AllSvc || svc == ResizeSvc {
		listener := service.NewUploadEventListener()
		listener.Listener()
		defer listener.Close()
	}

	err := r.Run(":80")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Run failed")
	}
}
