package stability_ai

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-resty/resty/v2"
	"github.com/guregu/dynamo"
	"os"
)

var (
	bucket         string
	tableName      string
	region         string
	cloudfrontHost string
	apiKey         string
	client         *resty.Client
)

func init() {
	var isExist bool

	apiKey, isExist = os.LookupEnv("STABILITY_KEY")
	if !isExist {
		panic("OS Environment does not exist [STABILITY_KEY]")
	}

	bucket, isExist = os.LookupEnv("S3_BUCKET")
	if !isExist {
		panic("OS Environment does not exist [S3_BUCKET]")
	}

	tableName, isExist = os.LookupEnv("DYNAMODB_TABLE")
	if !isExist {
		panic("OS Environment does not exist [DYNAMODB_TABLE]")
	}

	region, isExist = os.LookupEnv("DYNAMODB_REGION")
	if !isExist {
		panic("OS Environment does not exist [DYNAMODB_REGION]")
	}

	cloudfrontHost, isExist = os.LookupEnv("CLOUDFRONT_HOST")
	if !isExist {
		panic("OS Environment does not exist [CLOUDFRONT_HOST]")
	}

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table = db.Table(tableName)

	client = resty.New()
}
