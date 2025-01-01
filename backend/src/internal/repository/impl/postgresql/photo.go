package postgresql

import (
	"context"

	// "github.com/Masterminds/squirrel"
	// "github.com/jackc/pgx/v5"

	"backend/src/internal/model"
	"backend/src/internal/repository/interface"
	// "course/internal/model"
	"backend/src/internal/model/dto"
	// "course/internal/storage"
	// "course/pkg/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type photoMetaStorageImpl struct {
	db *pgxpool.Pool
}

func NewPhotoMetaStorage(db *pgxpool.Pool) _interface.IPhotoMetaStorage {
	return &photoMetaStorageImpl{db}
}

func (p *photoMetaStorageImpl) SaveKey(ctx context.Context, request *dto.CreatePhotoKeyRequest) error {
	query := `
		insert into photo(user_id, key) values ($1, $2)
	`

	_, err := p.db.Exec(
		ctx,
		query,
		request.UserId,
		request.Key,

	)

	if err != nil {
		return err
	}

	return nil
}

func (p *photoMetaStorageImpl) GetKey(ctx context.Context, request *dto.GetPhotoRequest) (*model.PhotoMeta, error) {
	query := `
		select id, user_id, key 
		from photo 
		where user_id = $1
	`

	photoMeta := new(model.PhotoMeta)

	err := p.db.QueryRow(
		ctx,
		query,
		request.UserId,
	).Scan(
		&photoMeta.Id,
		&photoMeta.UserId,
		&photoMeta.PhotoKey,

	)

	if err != nil {
		return nil, err
	}

	return photoMeta, nil
}

func (p *photoMetaStorageImpl) DeleteKey(ctx context.Context, request *dto.DeletePhotoRequest) error {
	query := `
		delete from photo where user_id = $1
	`

	_, err := p.db.Exec(
		ctx,
		query,
		request.UserId,
	)

	if err != nil {
		return err
	}

	return nil
}

// func (p *photoMetaStorageImpl) rowToModel(row pgx.Row) (*model.PhotoMeta, error) {
// 	var photoMeta model.PhotoMeta
// 	err := row.Scan(
// 		&photoMeta.ID,
// 		&photoMeta.DocumentID,
// 		&photoMeta.PhotoKey,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &photoMeta, nil
// }