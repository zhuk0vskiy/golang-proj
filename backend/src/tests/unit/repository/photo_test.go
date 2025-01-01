package repository_test

// import (
// 	"backend/src/internal/model/dto"
// 	"backend/src/internal/repository/interface/mocks"
// 	"backend/src/internal/service/impl"
// 	utils2 "backend/src/tests/utils"
// 	"context"
// 	"fmt"
// 	"github.com/ozontech/allure-go/pkg/framework/provider"
// 	"github.com/ozontech/allure-go/pkg/framework/suite"


// 	"context"
// 	"testing"

// 	"backend/src/internal/model/dto"
// 	"backend/src/internal/service/impl"
// 	"backend/src/pkg/logger"
// 	"backend/src/internal/repository/postgresql"
// 	"backend/src/internal/repository/mongodb"

// 	"github.com/stretchr/testify/assert"
// )

// type InstrumentalistRepositorySuite struct {
// 	suite.Suite

// 	mockRepo mocks.IInstrumentalistRepository
// }


    // Successfully save photo and its key for valid user ID and request data
// func TestCreatePhotoSuccess(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockPhotoStorage := mongodb.NewMockPhotoRepository(ctrl)
// 	mockLogger := postgresql.NewMockLogger(ctrl)

// 	photoService := &impl.PhotoService{
// 		photoStorage: mockPhotoStorage,
// 		logger: mockLogger,
// 	}

// 	ctx := context.Background()
// 	req := &dto.CreatePhotoRequest{
// 		UserId: 123,
// 		Data: []byte("test photo data"),
// 	}

// 	expectedKey := "photo-123.jpg"

// 	mockLogger.EXPECT().Infof("create photo by document ID %d", req.UserId)
// 	mockPhotoStorage.EXPECT().Save(ctx, req).Return(expectedKey, nil)
// 	mockPhotoStorage.EXPECT().SaveKey(ctx, &dto.CreatePhotoKeyRequest{
// 		UserId: req.UserId,
// 		Key: expectedKey,
// 	}).Return(nil)

// 	err := photoService.CreatePhoto(ctx, req)

// 	assert.NoError(t, err)
// }

	// Handle empty or invalid user ID in request
// func TestCreatePhotoInvalidUserId(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockPhotoStorage := mock_repository.NewMockPhotoRepository(ctrl)
// 	mockLogger := mock_logger.NewMockLogger(ctrl)

// 	photoService := &PhotoService{
// 		photoStorage: mockPhotoStorage,
// 		logger: mockLogger,
// 	}

// 	ctx := context.Background()
// 	req := &dto.CreatePhotoRequest{
// 		UserId: 0,
// 		Data: []byte("test photo data"),
// 	}

// 	mockLogger.EXPECT().Infof("create photo by document ID %d", req.UserId)
// 	mockPhotoStorage.EXPECT().Save(ctx, req).Return("", fmt.Errorf("invalid user ID"))
// 	mockLogger.EXPECT().Errorf("save photo: %s", "invalid user ID")

// 	err := photoService.CreatePhoto(ctx, req)

// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "save photo: invalid user ID")
// }

// package impl_test




// func TestCreatePhotoIntegration(t *testing.T) {
// 	// Setup the context
// 	ctx := context.Background()

// 	// Initialize the logger
// 	log := logger.New() // Assuming there's a New() function to create a logger instance

// 	// Initialize the photo storage
// 	photoStorage := interface.NewPhotoStorage() // Replace with actual initialization

// 	// Create an instance of PhotoService
// 	photoService := impl.PhotoService{
// 		logger:       log,
// 		photoStorage: photoStorage,
// 	}

// 	// Create a request
// 	request := &dto.CreatePhotoRequest{
// 		UserId: 123, // Example user ID
// 		// Add other necessary fields
// 	}

// 	// Call the CreatePhoto method
// 	err := photoService.CreatePhoto(ctx, request)

// 	// Assertions
// 	assert.NoError(t, err, "expected no error from CreatePhoto")
// 	// Additional assertions can be added here to verify the state of the system
// }Одеваемся
