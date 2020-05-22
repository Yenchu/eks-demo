package service

import (
	"eks-demo/awsapi"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"os"
)

func NewUploadEventListener() *UploadEventListener {

	queueURL := os.Getenv("UPLOAD_EVENT_QUEUE_URL")
	log.WithFields(log.Fields{"queueURL": queueURL}).Info("UPLOAD_EVENT_QUEUE_URL")

	sqsApi := awsapi.NewSqsAPI()

	imageSvc := NewImageService()

	return &UploadEventListener{
		queueUrl: queueURL,
		sqsApi:   sqsApi,
		imageSvc: imageSvc,
		close:    make(chan bool),
	}
}

type UploadEventListener struct {
	queueUrl string
	sqsApi   *awsapi.SqsAPI
	imageSvc *ImageService
	close    chan bool
}

func (svc *UploadEventListener) handleMessage(msg string) error {

	log.WithFields(log.Fields{"msg": msg}).Info("Receive message")

	event := &events.S3Event{}
	err := json.Unmarshal([]byte(msg), event)
	if err != nil {
		return err
	}

	for _, record := range event.Records {

		s3 := record.S3

		bucket := s3.Bucket.Name
		key := s3.Object.Key

		log.WithFields(log.Fields{"source": record.EventSource, "event": record.EventName,
			"bucket": bucket, "key": key}).Info("S3 event record")

		err := svc.imageSvc.ResizeImage(bucket, key)
		if err != nil {
			// TODO: need a better error handler
			log.WithFields(log.Fields{"error": err}).Error("handleMessage failed")
		}
	}
	return nil
}

func (svc *UploadEventListener) pollMessages() {

	for {
		select {
		case closed := <-svc.close:
			log.WithFields(log.Fields{"closed": closed}).Info("pollMessages stopped")
			return
		default:
		}

		err := svc.sqsApi.PollMessages(svc.queueUrl, svc.handleMessage)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("pollMessages failed")
			return
		}
	}
}

func (svc *UploadEventListener) Listener() {

	log.Info("Start polling messages")
	go svc.pollMessages()
}

func (svc *UploadEventListener) Close() {

	log.Info("Stop polling messages")
	svc.close <- true
}
