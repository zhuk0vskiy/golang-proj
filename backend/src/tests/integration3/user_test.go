package tests

import (
	"backend/src/internal/model"
	"backend/src/internal/model/dto"
	"backend/src/internal/repository/impl/postgresql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"testing"
)

func TestUserRepository_GetByLogin(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx     context.Context
		request *dto.GetUserByLoginRequest
	}
	tests := []struct {
		name string
		//fields   fields
		args     args
		wantUser *model.User
		wantErr  bool
	}{
		//{
		//	name: "test_pos_01",
		//	args: args{
		//		ctx:     context.Background(),
		//		request: &dto.GetUserByLoginRequest{Login: "test"},
		//	},
		//	wantUser: &model.User{
		//		Id:         1,
		//		Login:      "test",
		//		Password:   "test",
		//		Role:       "test",
		//		FirstName:  "test",
		//		SecondName: "test",
		//		ThirdName:  "test",
		//	},
		//	wantErr: false,
		//},
		{
			name: "test_neg_01",
			args: args{
				ctx:     context.Background(),
				request: &dto.GetUserByLoginRequest{Login: "test_neg"},
			},
			wantUser: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgresql.NewUserRepository(testDbInstance)
			gotUser, err := r.GetByLogin(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetByLogin() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
