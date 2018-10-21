package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pelletier/go-toml"
)

var preserveDirStructureBool bool

func main() {
	fmt.Println("Loading deploy.toml file...")
	preserveDirStructureBool = true
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dat, err := ioutil.ReadFile(dir + "/deploy.toml")
	if err != nil {
		log.Fatal("Could not read deploy.toml file", err)
	}

	config, _ := toml.Load(string(dat))

	fmt.Println("Building hugo site...")
	cmd := exec.Command("hugo", "-t", "uilite", "-D")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.hugo() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Get("aws.region").(string)),
		Credentials: credentials.NewStaticCredentials(config.Get("aws.keyid").(string), config.Get("aws.secretkey").(string), ""),
	})

	fmt.Println("Uploading to S3 - ", config.Get("aws.bucketname").(string))
	fmt.Println("=================================")
	uploadDirToS3(sess, config.Get("aws.bucketname").(string), "", dir+"/public")

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

func uploadDirToS3(sess *session.Session, bucketName string, bucketPrefix string, dirPath string) {
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
		uploadFileToS3(sess, bucketName, bucketPrefix, file, dirPath)
	}
}

func uploadFileToS3(sess *session.Session, bucketName string, bucketPrefix string, filePath string, dirPath string) {
	fmt.Println("upload " + filePath + " to S3")
	// An s3 service
	s3Svc := s3.New(sess)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file", file, err)
		os.Exit(1)
	}
	defer file.Close()
	var key string
	if preserveDirStructureBool {
		fileDirectory, _ := filepath.Abs(filePath)
		fileDirectory = strings.Replace(fileDirectory, dirPath+"/", "", 1)
		key = bucketPrefix + fileDirectory
	} else {
		key = bucketPrefix + path.Base(filePath)
	}
	// Upload the file to the s3 given bucket
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName), // Required
		Key:    aws.String(key),        // Required
		Body:   file,
	}
	_, err = s3Svc.PutObject(params)
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n",
			bucketName, key, err.Error())
		return
	}
}
