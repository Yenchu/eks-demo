package awsapi

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"time"
)

func NewSsmAPI() *SsmAPI {

	client := ssm.New(LoadAWSConfig())

	return &SsmAPI{
		client: client,
	}
}

type SsmAPI struct {
	client *ssm.Client
}

func (api *SsmAPI) GetParameter(name string) (string, error) {

	input := &ssm.GetParameterInput{
		Name: &name,
	}

	req := api.client.GetParameterRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return "", err
	}
	return *resp.Parameter.Value, nil
}

func (api *SsmAPI) GetDecryptedParameter(name string) (string, error) {

	decrypt := true
	input := &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &decrypt,
	}

	req := api.client.GetParameterRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return "", err
	}
	return *resp.Parameter.Value, nil
}

func (api *SsmAPI) GetParameters(names ...string) (map[string]string, error) {

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	input := &ssm.GetParametersInput{
		Names: names,
	}

	req := api.client.GetParametersRequest(input)

	resp, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, param := range resp.Parameters {
		result[*param.Name] = *param.Value
	}
	return result, nil
}
