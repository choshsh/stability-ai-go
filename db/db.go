package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
)

var db *dynamo.DB

func Init() {
	region, isExist := os.LookupEnv("DYNAMODB_REGION")
	if !isExist {
		panic("OS Environment does not exist [DYNAMODB_REGION]")
	}

	var config aws.Config

	switch os.Getenv("GO_ENV") {
	case "production":
		config = aws.Config{Region: aws.String(region)}
	default:
		config = aws.Config{
			Endpoint: aws.String("http://localhost:8000"),
			Region:   aws.String(region),
		}
	}

	db = dynamo.New(session.Must(session.NewSession()), &config)
}

func DB() *dynamo.DB {
	return db
}
