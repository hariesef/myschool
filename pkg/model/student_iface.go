package model

//go:generate mockgen -source=student_iface.go -package=mocks -destination=../../internal/mocks/student_iface.go

import (
	"context"
)

// Student model
type StudentCreationParam struct {
	Name   string
	Gender string
}

type StudentModel interface {
	GetUID() uint
	GetCreatedAt() int
	GetUpdatedAt() int
	GetName() string
	GetGender() string
}

type StudentRepo interface {
	Create(ctx context.Context, args StudentCreationParam) (StudentModel, error)
	Read(ctx context.Context, uid uint) (StudentModel, error)
	Delete(ctx context.Context, uid uint) (StudentModel, error)
	FindByName(ctx context.Context, studentName string) ([]StudentModel, error)
}
