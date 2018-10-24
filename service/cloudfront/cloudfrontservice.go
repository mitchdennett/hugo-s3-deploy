package cloudfront

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/mitchdennett/hugo-s3-deploy/service/s3"
)

type Distribution struct {
	session    *session.Session
	Bucket     *s3.S3Bucket
	Region     string
	AliasName  string
	DomainName *string
}

func NewDistribution(sess *session.Session) *Distribution {
	dist := new(Distribution)
	dist.session = sess
	return dist
}

func (dist *Distribution) SetAliasName(name string) {
	dist.AliasName = name
}

func (dist *Distribution) SetRegion(region string) {
	dist.Region = region
}

func (dist *Distribution) SetBucket(bucket *s3.S3Bucket) {
	dist.Bucket = bucket
}

func (dist *Distribution) CreateDistribution() {
	svc := cloudfront.New(dist.session)

	origin := &cloudfront.Origin{
		DomainName: aws.String(dist.Bucket.Name + ".s3-website-" + dist.Region + ".amazonaws.com"),
		Id:         aws.String("S3-Website-" + dist.Bucket.Name + ".s3-website-" + dist.Region + ".amazonaws.com"),
		CustomOriginConfig: &cloudfront.CustomOriginConfig{
			HTTPPort:             aws.Int64(80),
			HTTPSPort:            aws.Int64(443),
			OriginProtocolPolicy: aws.String("http-only"),
		},
	}
	origins := []*cloudfront.Origin{origin}

	input := &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			Aliases: &cloudfront.Aliases{
				Items:    aws.StringSlice([]string{dist.AliasName, "www." + dist.AliasName}),
				Quantity: aws.Int64(2),
			},
			CallerReference: aws.String(strconv.FormatInt(time.Now().UnixNano(), 10)),
			Comment:         aws.String("Cloudfront for " + dist.AliasName),
			Enabled:         aws.Bool(true),
			Origins: &cloudfront.Origins{
				Items:    origins,
				Quantity: aws.Int64(1),
			},
			DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
				ForwardedValues: &cloudfront.ForwardedValues{
					Cookies: &cloudfront.CookiePreference{
						Forward: aws.String("none"),
					},
					QueryString: aws.Bool(false),
				},
				MinTTL:               aws.Int64(0),
				TargetOriginId:       aws.String("S3-Website-" + dist.Bucket.Name + ".s3-website-" + dist.Region + ".amazonaws.com"),
				ViewerProtocolPolicy: aws.String("redirect-to-https"),
				TrustedSigners: &cloudfront.TrustedSigners{
					Enabled:  aws.Bool(false),
					Quantity: aws.Int64(0),
				},
			},
		},
	}

	result, err := svc.CreateDistribution(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeCNAMEAlreadyExists:
				fmt.Println(cloudfront.ErrCodeCNAMEAlreadyExists, aerr.Error())
			case cloudfront.ErrCodeDistributionAlreadyExists:
				fmt.Println(cloudfront.ErrCodeDistributionAlreadyExists, aerr.Error())
			case cloudfront.ErrCodeInvalidOrigin:
				fmt.Println(cloudfront.ErrCodeInvalidOrigin, aerr.Error())
			case cloudfront.ErrCodeInvalidOriginAccessIdentity:
				fmt.Println(cloudfront.ErrCodeInvalidOriginAccessIdentity, aerr.Error())
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			case cloudfront.ErrCodeTooManyTrustedSigners:
				fmt.Println(cloudfront.ErrCodeTooManyTrustedSigners, aerr.Error())
			case cloudfront.ErrCodeTrustedSignerDoesNotExist:
				fmt.Println(cloudfront.ErrCodeTrustedSignerDoesNotExist, aerr.Error())
			case cloudfront.ErrCodeInvalidViewerCertificate:
				fmt.Println(cloudfront.ErrCodeInvalidViewerCertificate, aerr.Error())
			case cloudfront.ErrCodeInvalidMinimumProtocolVersion:
				fmt.Println(cloudfront.ErrCodeInvalidMinimumProtocolVersion, aerr.Error())
			case cloudfront.ErrCodeMissingBody:
				fmt.Println(cloudfront.ErrCodeMissingBody, aerr.Error())
			case cloudfront.ErrCodeTooManyDistributionCNAMEs:
				fmt.Println(cloudfront.ErrCodeTooManyDistributionCNAMEs, aerr.Error())
			case cloudfront.ErrCodeTooManyDistributions:
				fmt.Println(cloudfront.ErrCodeTooManyDistributions, aerr.Error())
			case cloudfront.ErrCodeInvalidDefaultRootObject:
				fmt.Println(cloudfront.ErrCodeInvalidDefaultRootObject, aerr.Error())
			case cloudfront.ErrCodeInvalidRelativePath:
				fmt.Println(cloudfront.ErrCodeInvalidRelativePath, aerr.Error())
			case cloudfront.ErrCodeInvalidErrorCode:
				fmt.Println(cloudfront.ErrCodeInvalidErrorCode, aerr.Error())
			case cloudfront.ErrCodeInvalidResponseCode:
				fmt.Println(cloudfront.ErrCodeInvalidResponseCode, aerr.Error())
			case cloudfront.ErrCodeInvalidArgument:
				fmt.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			case cloudfront.ErrCodeInvalidRequiredProtocol:
				fmt.Println(cloudfront.ErrCodeInvalidRequiredProtocol, aerr.Error())
			case cloudfront.ErrCodeNoSuchOrigin:
				fmt.Println(cloudfront.ErrCodeNoSuchOrigin, aerr.Error())
			case cloudfront.ErrCodeTooManyOrigins:
				fmt.Println(cloudfront.ErrCodeTooManyOrigins, aerr.Error())
			case cloudfront.ErrCodeTooManyCacheBehaviors:
				fmt.Println(cloudfront.ErrCodeTooManyCacheBehaviors, aerr.Error())
			case cloudfront.ErrCodeTooManyCookieNamesInWhiteList:
				fmt.Println(cloudfront.ErrCodeTooManyCookieNamesInWhiteList, aerr.Error())
			case cloudfront.ErrCodeInvalidForwardCookies:
				fmt.Println(cloudfront.ErrCodeInvalidForwardCookies, aerr.Error())
			case cloudfront.ErrCodeTooManyHeadersInForwardedValues:
				fmt.Println(cloudfront.ErrCodeTooManyHeadersInForwardedValues, aerr.Error())
			case cloudfront.ErrCodeInvalidHeadersForS3Origin:
				fmt.Println(cloudfront.ErrCodeInvalidHeadersForS3Origin, aerr.Error())
			case cloudfront.ErrCodeInconsistentQuantities:
				fmt.Println(cloudfront.ErrCodeInconsistentQuantities, aerr.Error())
			case cloudfront.ErrCodeTooManyCertificates:
				fmt.Println(cloudfront.ErrCodeTooManyCertificates, aerr.Error())
			case cloudfront.ErrCodeInvalidLocationCode:
				fmt.Println(cloudfront.ErrCodeInvalidLocationCode, aerr.Error())
			case cloudfront.ErrCodeInvalidGeoRestrictionParameter:
				fmt.Println(cloudfront.ErrCodeInvalidGeoRestrictionParameter, aerr.Error())
			case cloudfront.ErrCodeInvalidProtocolSettings:
				fmt.Println(cloudfront.ErrCodeInvalidProtocolSettings, aerr.Error())
			case cloudfront.ErrCodeInvalidTTLOrder:
				fmt.Println(cloudfront.ErrCodeInvalidTTLOrder, aerr.Error())
			case cloudfront.ErrCodeInvalidWebACLId:
				fmt.Println(cloudfront.ErrCodeInvalidWebACLId, aerr.Error())
			case cloudfront.ErrCodeTooManyOriginCustomHeaders:
				fmt.Println(cloudfront.ErrCodeTooManyOriginCustomHeaders, aerr.Error())
			case cloudfront.ErrCodeTooManyQueryStringParameters:
				fmt.Println(cloudfront.ErrCodeTooManyQueryStringParameters, aerr.Error())
			case cloudfront.ErrCodeInvalidQueryStringParameters:
				fmt.Println(cloudfront.ErrCodeInvalidQueryStringParameters, aerr.Error())
			case cloudfront.ErrCodeTooManyDistributionsWithLambdaAssociations:
				fmt.Println(cloudfront.ErrCodeTooManyDistributionsWithLambdaAssociations, aerr.Error())
			case cloudfront.ErrCodeTooManyLambdaFunctionAssociations:
				fmt.Println(cloudfront.ErrCodeTooManyLambdaFunctionAssociations, aerr.Error())
			case cloudfront.ErrCodeInvalidLambdaFunctionAssociation:
				fmt.Println(cloudfront.ErrCodeInvalidLambdaFunctionAssociation, aerr.Error())
			case cloudfront.ErrCodeInvalidOriginReadTimeout:
				fmt.Println(cloudfront.ErrCodeInvalidOriginReadTimeout, aerr.Error())
			case cloudfront.ErrCodeInvalidOriginKeepaliveTimeout:
				fmt.Println(cloudfront.ErrCodeInvalidOriginKeepaliveTimeout, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}

		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		log.Fatal("Unable to create CloudFront Distribution")
	}

	dist.DomainName = result.Distribution.DomainName
}
