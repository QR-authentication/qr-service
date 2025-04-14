package service

import "context"

type DBRepo interface {
	StoreToken(ctx context.Context, token, uuid string) error
	TokenStatusIsScanned(ctx context.Context, token string) (bool, error)
	UpdateTokenStatusToExpired(ctx context.Context, token string) error
	UpdateTokenStatusToScanned(ctx context.Context, token string) error
	GetLatestAction(ctx context.Context, uuid string) (string, error)
	InsertAction(ctx context.Context, action, uuid string) error
}
