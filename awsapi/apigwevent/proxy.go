package apigwevent

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ParseRequest(req events.APIGatewayProxyRequest, data interface{}) error {

	if req.Body != "" {
		err := json.Unmarshal([]byte(req.Body), &data)
		if err != nil {
			return err
		}
	}
	return nil
}

func corsHeaders() map[string]string {

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type,Authorization,X-Amz-Date,X-Api-Key",
	}
	return headers
}

func Response(body interface{}, status int) (events.APIGatewayProxyResponse, error) {

	headers := corsHeaders()

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Parse JSON failed")
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}

	bodyStr := string(bodyBytes)

	if status < 200 || status >= 400 {
		log.WithFields(log.Fields{"status": status, "body": bodyStr}).Error("Handle request failed")
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: status,
		Body:       bodyStr,
	}, nil
}

func OK(body interface{}) (events.APIGatewayProxyResponse, error) {

	return Response(body, http.StatusOK)
}

func BadRequest(err error) (events.APIGatewayProxyResponse, error) {

	return Response(err.Error(), http.StatusBadRequest)
}

func InternalServerError(err error) (events.APIGatewayProxyResponse, error) {

	return Response(err.Error(), http.StatusInternalServerError)
}
