package model

type PhotoID int64

func ToPhotoID(id int64) *PhotoID {
	photoID := PhotoID(id)
	return &photoID
}

func (p *PhotoID) Int() int64 {
	return int64(*p)
}

type PhotoKey string

func ToPhotoKey(key string) *PhotoKey {
	photoKey := PhotoKey(key)
	return &photoKey
}

func (p *PhotoKey) String() string {
	return string(*p)
}

type PhotoMeta struct {
	Id         int64
	UserId int64
	PhotoKey   string
}

type Photo struct {
	Meta *PhotoMeta
	Data []byte
}