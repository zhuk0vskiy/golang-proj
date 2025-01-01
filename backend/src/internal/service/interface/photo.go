package _interface

import (
	"context"
	"backend/src/internal/model"
	"backend/src/internal/model/dto"
)
type PhotoService interface {
	CreatePhoto(ctx context.Context, request *dto.CreatePhotoRequest) error
	GetPhoto(ctx context.Context, request *dto.GetPhotoRequest) (*model.Photo, error)
	DeletePhoto(ctx context.Context, request *dto.DeletePhotoRequest) error
}