package advertisingentity

import "context"

type AdvertisingRepository interface {
	CreateAdvertising(ctx context.Context, Advertising *Advertising) (err error)
	UpdateAdvertising(ctx context.Context, Advertising *Advertising) (err error)
	DeleteAdvertising(ctx context.Context, id string) (err error)
	GetAdvertisingByID(ctx context.Context, id string) (Advertising *Advertising, err error)
	GetAllAdvertisements(ctx context.Context) ([]Advertising, error)
}
