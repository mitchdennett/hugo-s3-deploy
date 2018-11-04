package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitchdennett/hugo-s3-deploy/service/acm"
	"github.com/mitchdennett/hugo-s3-deploy/service/cloudfront"
	"github.com/mitchdennett/hugo-s3-deploy/service/route53"
	s3Service "github.com/mitchdennett/hugo-s3-deploy/service/s3"
	"github.com/pelletier/go-toml"
)

var preserveDirStructureBool bool
var bucketExists bool

var bucketName string
var domainName string
var hostedZoneId string
var region string

func main() {
	preserveDirStructureBool = true
	fmt.Println("Loading deploy.toml file...")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	config := loadConfigToml(dir)
	bucketName = config.Get("aws.bucketname").(string)
	domainName = config.Get("aws.domain").(string)
	hostedZoneId = config.Get("aws.hostedzoneid").(string)
	region = config.Get("aws.region").(string)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Get("aws.region").(string)),
		Credentials: credentials.NewStaticCredentials(config.Get("aws.keyid").(string), config.Get("aws.secretkey").(string), ""),
	})

	bucket := s3Service.NewBucket(sess)
	bucket.SetName(bucketName)
	bucket.SetRegion(region)

	cert := acm.NewCert(sess)
	cert.SetDomainName(domainName)
	cert.SetHostedZoneId(hostedZoneId)

	dist := cloudfront.NewDistribution(sess)
	dist.SetAliasName(domainName)
	dist.SetRegion(region)
	dist.SetBucket(bucket)

	bucketExists = bucket.CreateOrRetrieve()

	if !bucketExists {
		fmt.Println("Requesting Cert....")
		fmt.Println("=================================")
		cert.Request(sess)
		resourceRecord := cert.DescribeCertificate()

		fmt.Println("Inserting Cert DNS Verification")
		fmt.Println("=================================")
		route53.InsertNewRecord(sess, cert, resourceRecord)

		fmt.Println("Setting Bucket Policy....")
		fmt.Println("=================================")
		bucket.MakePublic()

		fmt.Println("Setting up bucket for hosting....")
		fmt.Println("=================================")
		bucket.EnableWebHosting()

		fmt.Println("Creating CloudFront Distribution....")
		fmt.Println("=================================")
		dist.CreateDistribution()

		fmt.Println("Adding CloudFront Domain To DNS....")
		fmt.Println("=================================")
		route53.ChangeHostedZoneRecord(dist.DomainName, sess, domainName, hostedZoneId)
	}

	fmt.Println("Building Hugo Site....")
	fmt.Println("=================================")
	buildHugoSite(dir)

	fmt.Println("Uploading to S3 - ", bucketName)
	fmt.Println("=================================")
	bucket.UploadDirectory("", dir+"/public")
}

func buildHugoSite(dir string) {
	cmd := exec.Command("hugo", "-t", "hugo-universal-theme")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.hugo() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func loadConfigToml(dir string) *toml.Tree {
	dat, err := ioutil.ReadFile(dir + "/deploy.toml")
	if err != nil {
		log.Fatal("Could not read deploy.toml file", err)
	}

	config, _ := toml.Load(string(dat))
	return config
}
