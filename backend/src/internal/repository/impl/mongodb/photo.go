package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "backend/src/internal/model"
	"backend/src/internal/model/dto"
	// repoInterface "backend/src/internal/repository/interface"
	"backend/src/pkg/mongodb"
)

type photoDataStorageImpl struct {
	*mongodb.MongoDB
}

func NewPhotoDataStorage(db *mongodb.MongoDB) *photoDataStorageImpl{
	return &photoDataStorageImpl{db}
}

func (p *photoDataStorageImpl) Save(ctx context.Context, request *dto.CreatePhotoRequest) (string, error) {
	userId := strconv.Itoa(int(request.UserId))
	uploadOpts := options.GridFSUpload().SetMetadata(bson.D{
		{
			Key:   userId,
			Value: len(request.Data),
		},
	})

	objectID, err := p.Bucket.UploadFromStream(
		fmt.Sprintf("%s_%d.jpg", userId, time.Now().Unix()),
		bytes.NewReader(request.Data),
		uploadOpts,
	)
	if err != nil {
		return "", err
	}

	return objectID.Hex(), nil
}

func (p *photoDataStorageImpl) Get(ctx context.Context, key string) ([]byte, error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	if _, err = p.Bucket.DownloadToStream(id, buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (p *photoDataStorageImpl) Delete(ctx context.Context, key string) error {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return err
	}

	return p.Bucket.Delete(id)
}