package route53

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/route53"
	acmService "github.com/mitchdennett/hugo-s3-deploy/service/acm"
)

func InsertNewRecord(sess *session.Session, cert *acmService.Certificate, resourceRecord1 *acm.ResourceRecord) {
	change1 := &route53.Change{
		Action: aws.String("UPSERT"),
		ResourceRecordSet: &route53.ResourceRecordSet{
			Name: resourceRecord1.Name,
			Type: resourceRecord1.Type,
			ResourceRecords: []*route53.ResourceRecord{
				&route53.ResourceRecord{Value: resourceRecord1.Value},
			},
			TTL: aws.Int64(60),
		},
	}

	changes := []*route53.Change{change1}

	r53 := route53.New(sess)

	_, changeErr := r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(cert.HostedZoneId),
	})

	if changeErr != nil {
		if aerr, ok := changeErr.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeNoSuchHealthCheck:
				fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
			case route53.ErrCodeInvalidChangeBatch:
				fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			case route53.ErrCodePriorRequestNotComplete:
				fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(changeErr.Error())
		}
		log.Fatal("Error adding CNAME records")
	}
}

func ChangeHostedZoneRecord(cloudFrontDomain *string, sess *session.Session, domainName string, hostedZoneId string) {
	change1 := &route53.Change{
		Action: aws.String("UPSERT"),
		ResourceRecordSet: &route53.ResourceRecordSet{
			Name: aws.String("www." + domainName),
			Type: aws.String("CNAME"),
			ResourceRecords: []*route53.ResourceRecord{
				&route53.ResourceRecord{Value: cloudFrontDomain},
			},
			TTL: aws.Int64(60),
		},
	}

	change2 := &route53.Change{
		Action: aws.String("UPSERT"),
		ResourceRecordSet: &route53.ResourceRecordSet{
			AliasTarget: &route53.AliasTarget{
				DNSName:              cloudFrontDomain,
				EvaluateTargetHealth: aws.Bool(false),
				HostedZoneId:         aws.String("Z2FDTNDATAQYW2"),
			},
			Name: aws.String(domainName),
			Type: aws.String("A"),
		},
	}

	changes := []*route53.Change{change1, change2}

	r53 := route53.New(sess)

	_, changeErr := r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(hostedZoneId),
	})

	if changeErr != nil {
		if aerr, ok := changeErr.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeNoSuchHealthCheck:
				fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
			case route53.ErrCodeInvalidChangeBatch:
				fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			case route53.ErrCodePriorRequestNotComplete:
				fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(changeErr.Error())
		}
		log.Fatal("Error adding CloudFront CNAME record")
	}
}
