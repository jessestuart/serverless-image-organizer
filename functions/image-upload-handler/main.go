package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
const TargetBucketName string = "jesse.pics"

// NetlifyWebhookURL is the netlify-provided webhook to trigger a rebuild of
// your site after image processing is complete.
// Can be left empty if that's not your jam.
var NetlifyWebhookURL = os.Getenv("JESSESIO_NETLIFY_WEBHOOK")

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

	fmt.Println("Lambda triggered by upload of S3 object with key: " + objectKey)

	result, err := s3Client.GetObject(getObjInput)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Cannot fetch object: %s", err), err
	}

	exifData, err := exif.Decode(result.Body)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Cannot decode EXIF data: %s", err), err
	}

	datetime, _ := exifData.DateTime()
	response := "Date created: " + datetime.Format(time.RFC3339)

	dateString := datetime.Format("2006-01-02")

	// Next, we'll copy the object from the source to the target S3 bucket, using
	// the datetime obtained from EXIF data to structure the directory tree.
	copyObjInput := &s3.CopyObjectInput{
		Bucket:     aws.String(TargetBucketName),
		CopySource: aws.String(fmt.Sprintf("/%s/%s", bucketName, objectKey)),
		Key:        aws.String(fmt.Sprintf("/%s/%s", dateString, objectKey)),
	}
	copyResult, err := s3Client.CopyObject(copyObjInput)

	if err != nil {
		// Copy S3 object operation failed.
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

	// Delete the source object to clean up after ourselves.
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	deleteObjResult, err := s3Client.DeleteObject(deleteObjectInput)
	fmt.Println(deleteObjResult)
	if err != nil {
		fmt.Println("Error deleting source object.")
		return fmt.Sprintf("Error: "), err
	}

	fmt.Println(fmt.Sprintf("S3 object copy completed successfully: %s", result))

	NotifyNetlify()

	return fmt.Sprintf(response), nil
}

// NotifyNetlify sends an HTTP POST to Netlify to trigger a build after image
// processing is complete.
func NotifyNetlify() (string, error) {
	url := NetlifyWebhookURL

	req, err := http.NewRequest("POST", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return fmt.Sprintf("Netlify response: %s", string(body)), nil
}

func main() {
	lambda.Start(HandleRequest)
}
