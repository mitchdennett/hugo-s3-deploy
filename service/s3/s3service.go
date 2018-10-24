package s3

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Bucket struct {
	Name    string
	Region  string
	session *session.Session
}

func NewBucket(session *session.Session) *S3Bucket {
	bucket := new(S3Bucket)
	bucket.session = session
	return bucket
}

func (bucket *S3Bucket) SetRegion(region string) {
	bucket.Region = region
}

func (bucket *S3Bucket) SetName(name string) {
	bucket.Name = name
}

func (bucket *S3Bucket) CreateOrRetrieve() bool {

	svc := s3.New(bucket.session)
	bucketExists := false

	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucket.Name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(bucket.Region),
		},
	}

	_, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				log.Fatal("Bucket: " + bucket.Name + " already exists. Please choose a different name")
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				bucketExists = true
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Fatal(err.Error())
		}
	} else {

	}

	return bucketExists
}

func (bucket *S3Bucket) MakePublic() {
	svc := s3.New(bucket.session)
	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket.Name),
		Policy: aws.String("{\"Version\":\"2008-10-17\",\"Statement\":[{\"Sid\":\"PublicReadGetObject\",\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"*\"},\"Action\":\"s3:GetObject\",\"Resource\":\"arn:aws:s3:::" + bucket.Name + "/*\"}]}"),
	}

	_, err := svc.PutBucketPolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Fatal(err.Error())
		}
		return
	}
}

func (bucket *S3Bucket) EnableWebHosting() {
	svc := s3.New(bucket.session)
	params := s3.PutBucketWebsiteInput{
		Bucket: aws.String(bucket.Name),
		WebsiteConfiguration: &s3.WebsiteConfiguration{
			IndexDocument: &s3.IndexDocument{
				Suffix: aws.String("index.html"),
			},
		},
	}

	_, err := svc.PutBucketWebsite(&params)
	if err != nil {
		log.Fatalf("Unable to set bucket %q website configuration, %v", bucket.Name, err)
	}
}

func isDirectory(path string) bool {
	fd, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}

func (bucket S3Bucket) UploadDirectory(bucketPrefix string, dirPath string) {
	fileList := []string{}
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		if isDirectory(path) {
			return nil
		} else {
			fileList = append(fileList, path)
			return nil
		}
	})

	for _, file := range fileList {
		bucket.UploadFile(bucketPrefix, file, dirPath)
	}
}

func (bucket *S3Bucket) UploadFile(bucketPrefix string, filePath string, dirPath string) {
	svc := s3.New(bucket.session)
	fmt.Println("upload " + filePath + " to S3")
	// An s3 service

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file", file, err)
		os.Exit(1)
	}
	defer file.Close()
	var key string
	fileDirectory, _ := filepath.Abs(filePath)
	fileDirectory = strings.Replace(fileDirectory, dirPath+"/", "", 1)
	key = bucketPrefix + fileDirectory
	// Upload the file to the s3 given bucket
	contentType := getFileContentType(filePath)
	params := &s3.PutObjectInput{
		Bucket:      aws.String(bucket.Name), // Required
		Key:         aws.String(key),         // Required
		Body:        file,
		ContentType: aws.String(contentType),
		Metadata: map[string]*string{
			"Content-Type": aws.String(contentType),
		},
	}
	_, err = svc.PutObject(params)
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n",
			bucket.Name, key, err.Error())
		return
	}
}

func getFileContentType(filePath string) string {
	out, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file", out, err)
		os.Exit(1)
	}
	defer out.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, readErr := out.Read(buffer)
	if readErr != nil {
		return ""
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return strings.Split(contentType, ";")[0]
}
