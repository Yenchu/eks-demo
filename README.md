## Image uploading and downloading EKS demo

This is a demo application for image uploading and downloading.
It uses S3 prs-signed URL to upload to S3 and CloudFront signed URL to download from CloudFront.

### Architecture Design

![Architecture Diagram](https://github.com/yenchu/eks-demo/raw/master/images/architecture.png)

### Build and Deployment

This project is deployed as a Kubernetes app, so you need to create an ECR image repository and EKS cluster.
To expose container applications to the internet, you need to create an Ingress Controller to route HTTP requests to pods.

##### Create Amazon ECR image repository

Please follow the AWS document
 [Getting Started with Amazon ECR using the AWS Management Console](https://docs.aws.amazon.com/AmazonECR/latest/userguide/getting-started-console.html)
 to create an image repository.
 
##### Create EKS cluster with Fargate

Please follow the AWS document 
 [Getting started with AWS Fargate on Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/fargate-getting-started.html)
 to create a Kubernetes cluster.

##### Create ALB ingress controller

Please follow the AWS document 
 [ALB Ingress Controller on Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html)
 to create an ALB ingress controller.

##### Build and Push docker image

To push images to ECR, you need to use an authorization token to authenticate docker with ECR:

```bash
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com
```

Build docker image:

```bash
docker build -t demo .
```

Tag image with ECR URI:

```bash
docker tag demo ${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/demo:latest
```

Push image to ECR repository:

```bash
docker push ${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/demo:latest
```

##### Create AWS resources

To create AWS resources(S3, CloudFormation and SQS) used by this demo:

```bash
aws cloudformation deploy --stack-name eks-demo-resources --template-file cloudformation/template.yaml --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND
```

##### Generate CloudFront Key Pair

Because signing CloudFront URL needs a key pair, you need to use root account to create one.

On AWS console, select `My Security Credentials` in menu bar, select `CloudFront key pairs`, click `Create New Key Pair`.
Then you need to store the generate key ID and private key in SSM parameter store. 

##### Store Key Pair in SSM Parameter Store

To store CloudFront key ID in SSM parameter store, go to `Systems Manager` console, select `Parameter Store`,
 save it with name `/applications/ServerlessDemo/CloudFront/KeyId` and type `String`.
 
To store private key in parameter store, save it with name `/applications/ServerlessDemo/CloudFront/PrivateKey` and type `SecureString`.

##### Use IRSA (IAM Roles for Service Accounts) to give pod permissions to access AWS resources

For container applications being able to access AWS resources, you need to use IRSA to give them permissions.

Create IAM policy:

```bash
aws iam create-policy --policy-name EksDemoPodIAMPolicy --policy-document k8s/pod-policy.json
```

Create service account and attache the previous created policy:

```bash
eksctl create iamserviceaccount \
    --region ap-northeast-1 \
    --name demo-isa \
    --cluster eks-demo \
    --attach-policy-arn arn:aws:iam::${AWS::AccountId}:policy/EksDemoPodIAMPolicy \
    --override-existing-serviceaccounts \
    --approve
```

Then you can use it as `serviceAccountName` in your pod spec.

##### Deploy demo app

Deploy a ConfigMap with S3, CloudFormation and SQS information:

```bash
kubectl apply -f k8s/configMap.yaml
```

Deploy upload-app for generating S3 pre-signed URL:

```bash
kubectl apply -f k8s/upload-app.yaml
```

Deploy download-app for generating CloudFront signed URL:

```bash
kubectl apply -f k8s/download-app.yaml
```

Deploy resize-app for resizing image:

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

### Test

You can use cURL to test these APIs.

##### Get S3 PreSigned URL for Upload

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

##### Upload to S3

To upload file to S3 using pre-signed URL, you need to pass back the headers you got from calling get-upload-url API.

``` 
curl -X PUT -H "content-type: image/jpg" -H "x-amz-meta-heigh: 1024" -H "x-amz-meta-width: 2048" {UPLOAD_URL}
 --data-binary '@{PATH_TO_IMAGE}.jpg'
```

##### Get CloudFront Signed URL for Download

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

##### Download from CloudFront

```
curl -X GET {DOWNLOAD_URL}
```
