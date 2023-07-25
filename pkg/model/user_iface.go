package model

//go:generate mockgen -source=user_iface.go -package=mocks -destination=../../internal/mocks/user_iface.go

import "context"

type UserCreationParam struct {
	Email             string
	EncryptedPassword string
}

type UserModel interface {
	GetID() string
	GetEmail() string
	GetEncryptedPassword() string
	IsActive() bool
}

type UserRepo interface {
	Create(ctx context.Context, args UserCreationParam) (UserModel, error)
	Read(ctx context.Context, email string) (UserModel, error)
	Deactivate(ctx context.Context, id string) (UserModel, error)
	FindActive(ctx context.Context) ([]UserModel, error)
}
