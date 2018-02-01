package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/rwcarlsen/goexif/exif"
)

// TargetBucketName is the bucket to which the source photos will be copied.
const TargetBucketName string = "js-photos"

// HandleRequest processes the incoming S3 PutObject event.
func HandleRequest(ctx context.Context, event events.S3Event) (string, error) {
	// Instantiate the S3 Clicnet.
	s3Client := s3.New(session.New())
	// Alias the S3 record for convenient referencing.
	s3Record := event.Records[0].S3

	// Extract the bucket name and object key from the lambda event...
	bucketName := s3Record.Bucket.Name
	objectKey := s3Record.Object.Key
	// And use these values to construct a `GetObjectInput` so we can retrieve
	// the actual image from S3.
	getObjInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	fmt.Println("Executing lambda for object: " + objectKey)

	result, err := s3Client.GetObject(getObjInput)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Error: %s", err), err
	}

	exifData, err := exif.Decode(result.Body)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Error: %s", err), err
	}

	datetime, _ := exifData.DateTime()
	response := "Date created: " + datetime.Format(time.RFC3339)

	dateString := datetime.Format("2006-01-02")
	// isoString := datetime.Format(time.RFC3339)

	copyObjInput := &s3.CopyObjectInput{
		Bucket:     aws.String(TargetBucketName),
		CopySource: aws.String(fmt.Sprintf("/%s/%s", bucketName, objectKey)),
		Key: aws.String(
			fmt.Sprintf("/%s/%s", dateString, objectKey),
		),
	}

	copyResult, err := s3Client.CopyObject(copyObjInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeObjectNotInActiveTierError:
				fmt.Println(s3.ErrCodeObjectNotInActiveTierError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return fmt.Sprintf(copyResult.GoString()), err
	}

	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	deleteObjResult, err := s3Client.DeleteObject(deleteObjectInput)
	fmt.Println(deleteObjResult)
	if err != nil {
		fmt.Println("Error deleting original object.")
		return fmt.Sprintf("Error: "), err
	}

	fmt.Println(result)

	return fmt.Sprintf(response), nil
}

func main() {
	lambda.Start(HandleRequest)
}
