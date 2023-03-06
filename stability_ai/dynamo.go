package stability_ai

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"stability-ai-go/db"
	"stability-ai-go/util"
	"strings"
	"time"
)

// SaveImage Store the generated image in S3 and related information in DynamoDB
func SaveImage(imageBytes []byte, request *StabilityApiPayload) (*Image, error) {
	// s3 upload
	id, _ := uuid.NewUUID()
	s3url, err := util.S3Upload(util.S3UploadInput{
		Bucket:        aws.String(bucket),
		ID:            aws.String(id.String()),
		RawBytes:      &imageBytes,
		ContentType:   aws.String("image/png"),
		FileExtension: aws.String(".png"),
	})
	if err != nil {
		return nil, err
	}

	// save
	table := db.DB().Table(tableName)

	image := Image{
		ImageID:     id,
		Time:        time.Now(),
		Url:         s3url,
		RequestInfo: request,
	}

	err = table.Put(image.SetKeywords()).If(ModelImagePutCondition).Run()
	if err != nil {
		log.Err(err).Caller().Msg("dynamodb put failed")
		return nil, err
	}
	return image.SetUrlPrefix(), nil
}

func FindImageById(id string) (*Image, error) {
	table := db.DB().Table(tableName)

	var images Images

	err := table.Scan().Filter("'ImageID' = ?", id).All(&images)
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, nil
	}

	images.SetUrlPrefix()
	return images[0], nil
}

type FindManyInput struct {
	Keyword string
}

func FindImageMany(input *FindManyInput) (Images, error) {
	table := db.DB().Table(tableName)

	var images Images
	var err error

	if len(strings.TrimSpace(input.Keyword)) > 0 {
		err = table.Scan().
			Filter("contains ($, ?)", "Keywords", input.Keyword).
			All(&images)
	} else {
		err = table.Scan().All(&images)
	}

	if err != nil {
		return nil, err
	}

	return images.SetUrlPrefix().SortTimeDesc(), nil
}

func DeleteImage(id string) error {
	table := db.DB().Table(tableName)

	image, err := FindImageById(id)
	if err != nil {
		return err
	}

	return table.Delete("ImageID", image.ImageID.String()).
		Range("Time", image.Time).
		Run()
}

func InitTable() error {
	d := db.DB()
	if _, err := d.Table(tableName).Describe().Run(); err != nil {
		log.Info().Msg("Create Table")
		return d.CreateTable(tableName, &Image{}).OnDemand(true).Wait()
	}
	return nil
}
