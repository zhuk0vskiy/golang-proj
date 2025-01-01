package _interface

import (
	"context"

	"backend/src/internal/model"
	"backend/src/internal/model/dto"
)

type IPhotoMetaStorage interface {
	SaveKey(ctx context.Context, request *dto.CreatePhotoKeyRequest) error
	GetKey(ctx context.Context, request *dto.GetPhotoRequest) (*model.PhotoMeta, error)
	DeleteKey(ctx context.Context, request *dto.DeletePhotoRequest) error
}

type IPhotoDataStorage interface {
	Save(ctx context.Context, request *dto.CreatePhotoRequest) (string, error)
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

//go:generate mockery --name=PhotoStorages
type IPhotoStorages interface {
	IPhotoDataStorage
	IPhotoMetaStorage
}

type PhotoStorage struct {
	IPhotoDataStorage
	IPhotoMetaStorage
}