package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "us-west-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(customResolver))
	cfg.Credentials = credentials.NewStaticCredentialsProvider("minio", "miniominio", "")

	if err != nil {
		log.Fatalln(err.Error())
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String("golang"),
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = client.PutPublicAccessBlock(context.Background(), &s3.PutPublicAccessBlockInput{
		Bucket: aws.String("golang"),
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       false,
			BlockPublicPolicy:     false,
			IgnorePublicAcls:      false,
			RestrictPublicBuckets: false,
		},
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = client.PutBucketWebsite(context.Background(), &s3.PutBucketWebsiteInput{
		Bucket: aws.String("golang"),
		WebsiteConfiguration: &types.WebsiteConfiguration{
			IndexDocument: &types.IndexDocument{
				Suffix: aws.String("index.html"),
			},
			ErrorDocument: &types.ErrorDocument{
				Key: aws.String("error.html"),
			},
			RedirectAllRequestsTo: &types.RedirectAllRequestsTo{
				HostName: aws.String(""),
			},
		},
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Body:                    strings.NewReader("test"),
		Bucket:                  aws.String("golang"),
		Key:                     aws.String("key"),
		WebsiteRedirectLocation: aws.String("https://www.google.com/"),
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	output, err := client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String("golang"),
		Key:    aws.String("key"),
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Printf("%+v", output)
}
