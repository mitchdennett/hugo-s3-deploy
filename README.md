# Hugo S3 Deploy CLI Tool
This Go program sets up everything you need to host a Hugo site on S3. First it creates a new S3 Buckets and enables it for Web Hosting. Then it requests a Certificate from Amazon Certificate Manager and then automatically inserts the CNAME record into your Route 53 hosted zone so the certificate can verify you own the domain. Next it creates a CloudFront distribution to sit in front of your S3 bucket and updates your Hosted Zone again to point to the CloudFront distribution. After that it builds your Hugo site and uploads it to your S3 bucket. The only thing it doesn't do is set your SSL certificate on the CloudFront Distribution. (It could do this automatically but the Certificate verification can take a while sometimes eg. > 30 mins.)

## Installation

To install, download this repository and then run the `go build -o hugo-s3-deploy` command. After that move the resulting binary to somewhere on your PATH so you can run it as a command line tool. 

### Configuration

You'll need to add a config.toml to the root of your Hugo site. 

```toml
[aws]
bucketname="NAME_OF_BUCKET"
keyid="AWS_KEY"
secretkey="AWS_SECRET"
region="AWS_REGION"
hostedzoneid="ROUTE 53 HOSTED ZONE ASSOCIATED WITH YOUR DOMAIN"
domain="EXAMPLE.COM (DO NOT INCLUDE WWW.)"

[hugo]
command="COMMAND TO BUILD HUGO"
```

### Running

Navigate to the root of your Hugo project and then run the following command

```bash
$ hugo-s3-deploy 
```

This will build your Hugo site and upload the files into the root of your S3 bucket and configure everything needed to have a Hugo site up and running on S3

### Notes

Please note this is very early stage software. I welcome any issues or contributions.

