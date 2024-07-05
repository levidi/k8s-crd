package aws

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"levi.com/bucket-operator/types"
)

func getSession(region string) *session.Session {

	endpoint, _ := os.LookupEnv("AWS_CONFIG_ENDPOINT")
	profile, _ := os.LookupEnv("AWS_PROFILE_NAME")
	id, _ := os.LookupEnv("AWS_CONFIG_ACCESS_KEY_ID")
	secret, _ := os.LookupEnv("AWS_CONFIG_SECRET_ACCESS_KEY")

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           aws.String(region),
			Endpoint:         aws.String(endpoint),
			Credentials:      credentials.NewStaticCredentials(id, secret, ""),
			LogLevel:         aws.LogLevel(aws.LogDebugWithHTTPBody),
			S3ForcePathStyle: aws.Bool(true),
		},
		Profile: profile,
	})
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	return sess
}

func CreateBucket(bucket *types.Bucket) {

	sess := getSession(bucket.Spec.Region)
	svc := s3.New(sess)

	fmt.Println("Creating bucket...")
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket.Spec.BucketName),
	})
	if err != nil {
		log.Fatalf("Error creating bucket: %v", err)
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket.Spec.BucketName),
	})
	if err != nil {
		log.Fatalf("Error waiting for bucket: %v", err)
	}

	fmt.Printf("Bucket %s created successfully.\n", bucket.Spec.BucketName)

	fmt.Println("Setting policy on the bucket...")
	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": "*",
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::%s/*"
            }
        ]
    }`, bucket.Spec.BucketName)

	_, err = svc.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket.Spec.BucketName),
		Policy: aws.String(policy),
	})
	if err != nil {
		log.Fatalf("Error setting policy on the bucket: %v", err)
	}

	fmt.Printf("Policy successfully set on the bucket. %s\n", bucket.Spec.BucketName)

}

func DeleteBucket(bucket *types.Bucket) {

	sess := getSession(bucket.Spec.Region)
	svc := s3.New(sess)

	fmt.Println("Deleting bucket...")
	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket.Spec.BucketName),
	})
	if err != nil {
		log.Fatalf("Error deleting bucket: %v", err)
	}

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket.Spec.BucketName),
	})
	if err != nil {
		log.Fatalf("Error waiting for bucket deletion: %v", err)
	}

	fmt.Printf("Bucket %s deleted successfully.\n", bucket.Spec.BucketName)
}
