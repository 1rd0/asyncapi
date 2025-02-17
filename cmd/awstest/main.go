package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-1"))

	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://s3.localhost.localstack.cloud:4566")
	})
	out, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err)
	}
	for _, bucket := range out.Buckets {
		fmt.Println(*bucket.Name)
	}
}
