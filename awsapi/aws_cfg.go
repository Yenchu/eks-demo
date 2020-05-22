package awsapi

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"os"
)

func LoadAWSConfig() aws.Config {

	region := os.Getenv("REGION")
	if region == "" {
		region = "ap-northeast-1" // default
	}

	awsCfg, err := external.LoadDefaultAWSConfig(external.WithDefaultRegion(region))
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	//fmt.Printf("aws region: %v\n", awsCfg.Region)
	return awsCfg
}
