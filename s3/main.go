package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,  func(lo *config.LoadOptions) error {
		lo.Region = "us-east-1"
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Println(err)
	}

	for _, bucket := range output.Buckets {
		fmt.Println(*bucket.Name)
	}
}
