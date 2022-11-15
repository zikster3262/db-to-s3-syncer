package awss3

import (
	"concurrency/src/utils"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func CreateBucket(s3Client *s3.Client, name string) error {

	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                    aws.String(name),
		ACL:                       types.BucketCannedACLPrivate,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraintUsWest2},
	})

	if err != nil {
		utils.FailOnError(err, err.Error())
	}

	return err
}
