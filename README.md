# Image Uploading and Downloading EKS Demo

This is a demo application for image uploading, downloading, and resizing on EKS.
It uses S3 prs-signed URL to upload to S3 and CloudFront signed URL to download from CloudFront.
After uploaded to S3, S3 will send an event notification to SQS, a container app receives this message and resize the image. 

## Architecture Design

![Architecture Diagram](https://github.com/yenchu/eks-demo/raw/master/images/architecture.png)

## Create AWS Resources

This project uses S3, CloudFormation, SQS, and parameter store, to create these services, you can use the template file `cloudformation/template.yaml`:

```bash
aws cloudformation deploy --stack-name eks-demo-resources --template-file cloudformation/template.yaml --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND
```

### Generate CloudFront Key Pair

Because signing CloudFront URL needs a key pair, you need to use root account to create one.

On AWS console, select `My Security Credentials` in menu bar, select `CloudFront key pairs`, click `Create New Key Pair`.
Then you need to store the generate key ID and private key in SSM parameter store. 

### Store Key Pair in SSM Parameter Store

To store CloudFront key ID in SSM parameter store, go to `Systems Manager` console, select `Parameter Store`,
 save it with name `/applications/ServerlessDemo/CloudFront/KeyId` and type `String`.
 
To store private key in parameter store, save it with name `/applications/ServerlessDemo/CloudFront/PrivateKey` and type `SecureString`.

## Setup Infrastructure

To deploy a Kubernetes app, you need to create an ECR image repository and EKS cluster.
If you have a microservice architecture, using an Ingress Controller is a better way
 to expose container applications to the internet and route HTTP requests to pods.

### Create Amazon ECR Image Repository

Please follow the AWS document
 [Getting Started with Amazon ECR using the AWS Management Console](https://docs.aws.amazon.com/AmazonECR/latest/userguide/getting-started-console.html)
 to create an image repository.
 
### Create EKS Cluster with Fargate

Please follow the AWS document 
 [Getting started with AWS Fargate on Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/fargate-getting-started.html)
 to create a Kubernetes cluster.

### Create ALB Ingress Controller

Please follow the AWS document 
 [ALB Ingress Controller on Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html)
 to create an ALB ingress controller. 
 You need to update the ingress-controller configuration with your cluster-name, aws-vpc-id and aws-region.

## Build and Deployment

### Build and Push Docker Image

To push images to ECR, you need to use an authorization token to authenticate docker with ECR:

```bash
aws ecr get-login-password --region ${AWS_Region} | docker login --username AWS --password-stdin ${AWS_AccountId}.dkr.ecr.${AWS_Region}.amazonaws.com
```

Build docker image:

```bash
docker build -t demo .
```

Tag image with ECR URI:

```bash
docker tag demo ${AWS_AccountId}.dkr.ecr.${AWS_Region}.amazonaws.com/demo:latest
```

Push image to ECR repository:

```bash
docker push ${AWS_AccountId}.dkr.ecr.${AWS_Region}.amazonaws.com/demo:latest
```

### Setup IRSA (IAM Roles for Service Accounts)

For container applications being able to access AWS resources, you need to use IRSA to give them permissions.

Create IAM policy (You need to update the json file with your AWS region, account ID, bucket name and queue name):

```bash
aws iam create-policy --policy-name EksDemoPodIAMPolicy --policy-document k8s/pod-policy.json
```

Create service account and attache the previous created policy:

```bash
eksctl create iamserviceaccount \
    --region ${AWS_Region} \
    --name demo-isa \
    --cluster ${EKS_CLUSTER_NAME} \
    --attach-policy-arn arn:aws:iam::${AW_:AccountId}:policy/EksDemoPodIAMPolicy \
    --override-existing-serviceaccounts \
    --approve
```

Then you can use it as `serviceAccountName` in your pod spec.

### Deploy Demo App

Deploy a ConfigMap with S3, CloudFormation and SQS information
 (You need to update the yaml file with your bucket name, CloudFront domain name and queue URL):

```bash
kubectl apply -f k8s/configMap.yaml
```

Deploy upload-app for generating S3 pre-signed URL
 (You need to update the yaml file with your image URI):

```bash
kubectl apply -f k8s/upload-app.yaml
```

Deploy download-app for generating CloudFront signed URL
 (You need to update the yaml file with your image URI):

```bash
kubectl apply -f k8s/download-app.yaml
```

Deploy resize-app for resizing image
 (You need to update the yaml file with your image URI):

```bash
kubectl apply -f k8s/resize-app.yaml
```

Deploy ingress to define routing rules:

```bash
kubectl apply -f k8s/ingress.yaml
```

Use the following command to find ALB endpoint:

```bash
kubectl get ingress
```

## Test

You can use cURL to test these APIs.

### Get S3 PreSigned URL for Upload

To get a S3 pre-signed URL for uploading, you need to provide a file name.
If you want to resize the image, you have to provide content type, width and height.

```
curl -X POST -H "Content-Type: application/json" http://{ALB_ENDPOINT}/upload/get-signed-url
 -d '{"file": "{FILE_NAME}", "contentType": "image/jpg", "width": 2048, "height": 1024}'
```

You will get the following response, the `url` is a S3 pre-signed URL you can use to upload file to S3. 
Please note the response headers need to be passed back to S3 when uploading file: 

```json
{
  "headers": {
    "content-type": "image/jpg",
    "x-amz-meta-height": "1024",
    "x-amz-meta-width": "2048"
  },
  "url": "{UPLOAD_URL}"
}
```

### Upload to S3

To upload file to S3 using pre-signed URL, you need to pass back the headers you got from calling get-upload-url API.

``` 
curl -X PUT -H "content-type: image/jpg" -H "x-amz-meta-heigh: 1024" -H "x-amz-meta-width: 2048" {UPLOAD_URL}
 --data-binary '@{PATH_TO_IMAGE}.jpg'
```

### Get CloudFront Signed URL for Download

To get a CloudFront signed URL for downloading, you need to provide the file name you want to download.

```
curl -X POST -H "Content-Type: application/json" http://{ALB_ENDPOINT}/download/get-signed-url
 -d '{"file": "{FILE_NAME}"'
```

You will get the following response, the `url` is a CloudFront signed URL you can use to download file from CloudFront. 

```json
{
  "url": "{DOWNLOAD_URL}"
}
```

### Download from CloudFront

```
curl -X GET {DOWNLOAD_URL}
```
