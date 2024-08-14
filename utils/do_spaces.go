package utils

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	// "time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)



func SaveLogsLocally(file multipart.File) {
	// header *multipart.File --> input variable
	//fileName := header.Filename

	out, err := os.Create("my_log.log")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
}



func UploadSpacesBucket(spaces string, filePath string, file multipart.File) {

	// // Current Time
	// loc, _ := time.LoadLocation("UTC")
	// utcTime := time.Now().In(loc)
	// utcTimeString := utcTime.String()



	// load env file
	LoadEnv()

	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	object := s3.PutObjectInput{
		Bucket: aws.String(spaces),
		Key:    aws.String(filePath),
		Body:   file,
		ACL:    aws.String("public"),
		// Metadata: map[string]*string{
		// 	"x-amz-meta-accound-id":   aws.String(accountId),
		// 	"x-amz-meta-spider-name":  aws.String(spiderName),
		// 	"x-amz-meta-job-name":     aws.String(jobName),
		// 	"x-amz-meta-job-group-id": aws.String(jobGroupId),
		// 	"x-amz-meta-utc-time":     aws.String(utcTimeString),
		// },
	}

	_, err := s3Client.PutObject(&object)
	if err != nil {
		log.Println(err.Error())
	}

}
