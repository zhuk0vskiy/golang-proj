package mongodb

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/internal/model/dto"
	"backend/src/pkg/mongodb"
)

const connURI = "mongodb://localhost:30001/"

var testMongoDB *mongodb.MongoDB

func TestMain(m *testing.M) {
	testMongoDB, _ = mongodb.New(connURI, "tests", "test")

	os.Exit(m.Run())
}

func Test_photoDataStorageImpl_Save(t *testing.T) {
	photoDataStorage := NewPhotoDataStorage(testMongoDB)

	key, err := photoDataStorage.Save(context.TODO(), &dto.CreatePhotoRequest{
		UserId: 1,
		Data:       []byte{'o'},
	})
	require.NoError(t, err)
	require.NotEmpty(t, key)
}

func Test_photoDataStorageImpl_Get(t *testing.T) {
	photoDataStorage := NewPhotoDataStorage(testMongoDB)

	request := &dto.CreatePhotoRequest{
		UserId: 1,
		Data:       []byte{'o'},
	}

	key, err := photoDataStorage.Save(context.TODO(), request)
	require.NoError(t, err)
	require.NotEmpty(t, key)
	fmt.Printf(key)

	data, err := photoDataStorage.Get(context.TODO(), key)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	require.Equal(t, data, request.Data)
}

func Test_photoDataStorageImpl_Delete(t *testing.T) {
	photoDataStorage := NewPhotoDataStorage(testMongoDB)

	// request := &dto.DeletePhotoRequest{
	// 	UserId: 1,
	// }

	err := photoDataStorage.Delete(context.TODO(), "6775207b7310e5112952c856")
	require.NoError(t, err)


	// require.NotEmpty(t, key)

	// err = photoDataStorage.Delete(context.TODO(), key)
	// require.NoError(t, err)

	// data, err := photoDataStorage.Get(context.TODO(), key)
	// require.Error(t, err)
	// require.Empty(t, data)
}