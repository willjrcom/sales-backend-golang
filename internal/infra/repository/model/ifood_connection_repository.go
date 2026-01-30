package model

import "context"

type IfoodConnectionRepository interface {
	Upsert(ctx context.Context, conn *IfoodConnection) error
	GetBySchema(ctx context.Context, schema string) (*IfoodConnection, error)
}
