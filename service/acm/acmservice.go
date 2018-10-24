package acm

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
)

type Certificate struct {
	Id           *string
	DomainName   string
	HostedZoneId string
	session      *session.Session
}

func NewCert(sess *session.Session) *Certificate {
	cert := new(Certificate)
	cert.session = sess
	return cert
}

func (cert *Certificate) SetDomainName(domainName string) {
	cert.DomainName = domainName
}

func (cert *Certificate) SetHostedZoneId(zoneId string) {
	cert.HostedZoneId = zoneId
}

func (cert *Certificate) setId(id *string) {
	cert.Id = id
}

func (cert *Certificate) Request(sess *session.Session) {
	svc := acm.New(sess, aws.NewConfig().WithRegion("us-east-1"))
	result, err := svc.RequestCertificate(&acm.RequestCertificateInput{
		DomainName:              aws.String("*." + cert.DomainName),
		ValidationMethod:        aws.String("DNS"),
		SubjectAlternativeNames: aws.StringSlice([]string{cert.DomainName}),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case acm.ErrCodeLimitExceededException:
				fmt.Println(acm.ErrCodeLimitExceededException, aerr.Error())
			case acm.ErrCodeInvalidDomainValidationOptionsException:
				fmt.Println(acm.ErrCodeInvalidDomainValidationOptionsException, aerr.Error())
			case acm.ErrCodeInvalidArnException:
				fmt.Println(acm.ErrCodeInvalidArnException, aerr.Error())
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Fatal(err.Error())
		}
		log.Fatal("Error Requesting Cert...")
	} else {
		//Call new method here
		time.Sleep(8 * time.Second)

		cert.setId(result.CertificateArn)
	}
}

func (cert *Certificate) DescribeCertificate() *acm.ResourceRecord {
	svc := acm.New(cert.session, aws.NewConfig().WithRegion("us-east-1"))
	result, err := svc.DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: cert.Id,
	})

	if err != nil {
		log.Fatal("Failed Describing Cert")
	}

	resourceRecord1 := result.Certificate.DomainValidationOptions[0].ResourceRecord

	if resourceRecord1 == nil {
		log.Fatal("Resource Record Doesn't exists.")
	}

	return resourceRecord1
}
