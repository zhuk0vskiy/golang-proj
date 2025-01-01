package impl

import (
	"context"
	"fmt"

	"backend/src/internal/model"
	"backend/src/internal/model/dto"
	repoInterface "backend/src/internal/repository/interface"
	"backend/src/pkg/logger"
	// serviceInterface "backend/src/internal/service/interface"
)

type PhotoService struct {
	logger       logger.Interface
	photoStorage repoInterface.IPhotoStorages
}

func NewPhotoService(logger logger.Interface, photoStorage repoInterface.IPhotoStorages) *PhotoService {
	return &PhotoService{
		logger:       logger,
		photoStorage: photoStorage,
	}
}

func (p *PhotoService) CreatePhoto(ctx context.Context, request *dto.CreatePhotoRequest) error {
	// FIXME: Crop face from document
	p.logger.Infof("create photo by document ID %d", request.UserId)

	key, err := p.photoStorage.Save(ctx, request)
	if err != nil {
		p.logger.Errorf("save photo: %s", err.Error())
		return fmt.Errorf("save photo: %w", err)
	}

	err = p.photoStorage.SaveKey(ctx, &dto.CreatePhotoKeyRequest{
		UserId: request.UserId,
		Key:    key,
	})
	if err != nil {
		p.logger.Errorf("save photo key: %s", err.Error())
		return fmt.Errorf("save photo key: %w", err)
	}

	return nil
}

func (p *PhotoService) GetPhoto(ctx context.Context, request *dto.GetPhotoRequest) (*model.Photo, error) {
	p.logger.Infof("get photo by document ID %d", request.UserId)

	meta, err := p.photoStorage.GetKey(ctx, request)
	if err != nil {
		p.logger.Errorf("get photo key: %s", err.Error())
		return nil, fmt.Errorf("get photo key: %w", err)
	}

	data, err := p.photoStorage.Get(ctx, meta.PhotoKey)
	if err != nil {
		p.logger.Errorf("get photo by key: %s", err.Error())
		return nil, fmt.Errorf("get photo by key: %w", err)
	}

	return &model.Photo{
		Meta: meta,
		Data: data,
	}, nil
}

func (p *PhotoService) DeletePhoto(ctx context.Context, request *dto.DeletePhotoRequest) error {
	p.logger.Infof("delete photo by document ID %d", request.UserId)

	meta, err := p.photoStorage.GetKey(ctx, &dto.GetPhotoRequest{UserId: request.UserId})
	if err != nil {
		p.logger.Errorf("get photo key: %s", err.Error())
		return fmt.Errorf("get photo key: %w", err)
	}

	err = p.photoStorage.Delete(ctx, meta.PhotoKey)
	if err != nil {
		p.logger.Errorf("delete photo by key: %s", err.Error())
		return fmt.Errorf("delete photo by key: %w", err)
	}

	err = p.photoStorage.DeleteKey(ctx, request)
	if err != nil {
		p.logger.Errorf("delete photo key: %s", err.Error())
		return fmt.Errorf("delete photo key: %w", err)
	}

	return nil
}
