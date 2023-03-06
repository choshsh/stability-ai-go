package stability_ai

import (
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"stability-ai-go/db"
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
	var err error

	if os.Getenv("GO_ENV") != "production" {
		projectName := regexp.MustCompile(`^(.*` + "stability-ai-go" + `)`)
		currentWorkDirectory, _ := os.Getwd()
		rootPath := projectName.Find([]byte(currentWorkDirectory))

		err = godotenv.Load(string(rootPath) + `/.env`)
		if err != nil {
			panic(err)
		}
	}

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

	cloudfrontHost, isExist = os.LookupEnv("CLOUDFRONT_HOST")
	if !isExist {
		panic("OS Environment does not exist [CLOUDFRONT_HOST]")
	}

	// Initialize HTTP Client
	client = resty.New()

	// Initialize DB
	db.Init()
	if _, err = db.DB().Table(tableName).Describe().Run(); err != nil {
		log.Info().Msg("Create Table")

		err = db.DB().CreateTable(tableName, &Image{}).OnDemand(true).Wait()
		if err != nil {
			panic(err)
		}
	}
}
