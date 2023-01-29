package stability_ai

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/rs/zerolog/log"
	"stability-ai-go/common"
	"strings"
	"time"
)

var table dynamo.Table

// SaveImage 생성된 이미지를 S3에 저장하고 관련 정보를 DynamoDB에 저장
func SaveImage(imageBytes []byte, request *StabilityApiPayload) (*Image, error) {
	// s3 upload
	id, _ := uuid.NewUUID()
	s3url, err := common.S3Upload(common.S3UploadInput{
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
	image := Image{
		ImageID:     id,
		Time:        time.Now(),
		Url:         s3url,
		RequestInfo: request,
	}
	image.SetKeywords()

	err = table.Put(image).If("attribute_not_exists(ImageID)").Run()
	if err != nil {
		log.Error().Err(err).Msg("dynamodb 저장 실패")
		return nil, err
	}

	image.SetUrlPrefix()
	return &image, nil
}

func FindImageById(id string) (*Image, error) {
	var images []*Image

	err := table.Scan().
		Filter("'ImageID' = ?", id).
		All(&images)

	if err != nil {
		return nil, err
	}

	if len(images) < 1 {
		return nil, nil
	}

	images[0].SetUrlPrefix()
	return images[0], nil
}

type FindManyInput struct {
	Keyword string
}

func FindImageMany(input *FindManyInput) ([]*Image, error) {
	var images []*Image
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

	for _, image := range images {
		image.SetUrlPrefix()
	}

	return images, nil
}

func DeleteImage(id string) error {
	image, err := FindImageById(id)
	if err != nil {
		return err
	}

	return table.Delete("ImageID", image.ImageID.String()).
		Range("Time", image.Time).
		Run()
}
