package model

//go:generate mockgen -source=token_iface.go -package=mocks -destination=../../internal/mocks/token_iface.go

import "context"

// Student model
type TokenCreationParam struct {
	Token  string
	UserID string
	Email  string
	Expiry int
}

type AuthTokenModel interface {
	GetID() string
	GetToken() string
	GetUserID() string
	GetEmail() string
	GetExpiry() int
}

type AuthTokenRepo interface {
	Create(ctx context.Context, args TokenCreationParam) (AuthTokenModel, error)
	Find(ctx context.Context, token string) (AuthTokenModel, error)
	Delete(ctx context.Context, token string) error
}
