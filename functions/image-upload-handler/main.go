package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/kelseyhightower/envconfig"
	"github.com/rwcarlsen/goexif/exif"
)

// EnvVars are user-provided variables that influence the Lambda function's
// execution.
type EnvVars struct {
	// TargetBucketName is the bucket to which the source photos will be copied.
	TargetBucketName       string `required:"true" split_words:"true"`
	WillDeleteSourceObject bool   `default:"true" split_words:"true"`
	CallbackWebhookURL     string `required:"false" split_words:"true"`
}

func parseEnvVars() (EnvVars, error) {
	var envVars EnvVars
	err := envconfig.Process("SIO", &envVars)
	if err != nil {
		return envVars, err
	}
	return envVars, nil
}

func invokeCallbackWebhook(callbackWebhookURL string) {
	response, err := http.Post(callbackWebhookURL, "", bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Println("Warning: a callback webhook was defined, but an error " +
			"occurred trying to invoke it.")
	}
	defer response.Body.Close()
}

// HandleRequest processes the incoming S3 PutObject event.
func HandleRequest(ctx context.Context, event events.S3Event) (string, error) {
	// Parse the environment variables.
	envVars, err := parseEnvVars()
	if err != nil {
		return fmt.Sprintf("Error: %s", err), err
	}
	fmt.Printf("%s\n%s\n", envVars.TargetBucketName, envVars.CallbackWebhookURL)

	// Instantiate the S3 Clicnet.
	s3Client := s3.New(session.New())
	// Alias the S3 record for convenient referencing.
	s3Record := event.Records[0].S3

	// Extract the bucket name and object key from the lambda event...
	bucketName := s3Record.Bucket.Name
	objectKey := s3Record.Object.Key
	objectPath := "/" + bucketName + "/" + objectKey

	fmt.Printf("%s\t%s\n", bucketName, objectKey)
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

	targetBucketName := envVars.TargetBucketName
	copyObjInput := &s3.CopyObjectInput{
		Bucket:     aws.String(targetBucketName),
		CopySource: aws.String(objectPath),
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

	if envVars.WillDeleteSourceObject {
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
	}

	webhookURL := envVars.CallbackWebhookURL
	if envVars.CallbackWebhookURL != "" {
		invokeCallbackWebhook(webhookURL)
	}

	return fmt.Sprintf(response), nil
}

func main() {
	lambda.Start(HandleRequest)
}
