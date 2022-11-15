package awss3

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func SetS3Config() *s3.Client {

	minioURL := os.Getenv("MINIO_URL")
	region := "us-east-1"
	user := os.Getenv("MINIO_ACCESS_KEY")
	pass := os.Getenv("MINIO_SECRET_KEY")

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			SigningRegion:     region,
			URL:               minioURL,
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:                      region,
		Credentials:                 credentials.NewStaticCredentialsProvider(user, pass, ""),
		EndpointResolverWithOptions: resolver,
	}

	return s3.NewFromConfig(cfg)

}
