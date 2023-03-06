package stability_ai

import (
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
	"stability-ai-go/db"
	"testing"
)

func TestDynamo(t *testing.T) {
	db.Init()

	err := InitTable()
	assert.NoError(t, err)
}

func TestDynamoGet(t *testing.T) {
	db.Init()
	table := db.DB().Table(tableName)

	var images []Image

	err := table.Scan().All(&images)
	assert.NoError(t, err)

	pp.Println("Fist:", images[0])
	pp.Println("Size:", len(images))
}

func TestDeleteAll(t *testing.T) {
	db.Init()
	table := db.DB().Table(tableName)

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
