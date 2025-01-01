package dto

type CreatePhotoRequest struct {
	UserId int64
	Data   []byte
}

type CreatePhotoKeyRequest struct {
	UserId int64
	Key    string
}

type GetPhotoRequest struct {
	UserId int64
}

type UpdatePhotoRequest struct {
	UserId int64
	Data   []byte
}

type UpdatePhotoKeyRequest struct {
	UserId int64
	Key    string
}

type DeletePhotoRequest struct {
	UserId int64
}
