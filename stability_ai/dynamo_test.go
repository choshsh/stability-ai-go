package stability_ai

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDynamo(t *testing.T) {
	sess, err := session.NewSessionWithOptions(session.Options{Profile: "verify"})
	assert.NoError(t, err)
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-2")})

	err = db.CreateTable("stability-ai", &Image{}).
		OnDemand(true).
		Run()

	assert.NoError(t, err)
}

func TestDynamoGet(t *testing.T) {
	sess, err := session.NewSessionWithOptions(session.Options{Profile: "verify"})
	assert.NoError(t, err)
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-2")})

	table = db.Table("stability-ai")

	var images []Image

	err = table.Scan().All(&images)
	assert.NoError(t, err)

	pp.Println("Size:", len(images))
}

func TestDeleteAll(t *testing.T) {
	sess, err := session.NewSessionWithOptions(session.Options{Profile: "verify"})
	assert.NoError(t, err)
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-2")})

	table = db.Table("stability-ai")

	images, err := FindImageMany(nil)
	assert.NoError(t, err)

	for _, image := range images {
		err = table.Delete("ImageID", image.ImageID.String()).
			Range("Time", image.Time).
			Run()
		assert.NoError(t, err)
	}

	images, err = FindImageMany(nil)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 0)

}
