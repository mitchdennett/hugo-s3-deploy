# Hugo S3 Deploy CLI Tool
This Go program will build your hugo site and then deploy to an S3 bucket.

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
```

### Running

Navigate to the root of your Hugo project and then run the following command

```bash
$ hugo-s3-deploy 
```

This will build your Hugo site and upload the files into the root of your S3 bucket.

### Notes

Please note this is very early stage software. I welcome any issues or contributions.

