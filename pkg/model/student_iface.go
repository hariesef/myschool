package model

//go:generate mockgen -source=student_iface.go -package=mocks -destination=../../internal/mocks/student_iface.go

import (
	"context"
)

// We are using this model to demonstrate SQL implementation of a table which contains basically only Name and Gender
// Uses: SQLite, gorm, sqlmock
type StudentModel interface {
	GetUID() uint      //gorm internal feature
	GetCreatedAt() int //gorm internal feature
	GetUpdatedAt() int //gorm internal feature
	GetDeletedAt() int //gorm internal feature
	GetName() string
	GetGender() string
}

// For create command specific
type StudentCreationParam struct {
	Name   string
	Gender string
}

// This is interface to be used by other layer, like service or controller
type StudentRepo interface {
	Create(ctx context.Context, args StudentCreationParam) (StudentModel, error)
	Read(ctx context.Context, uid uint) (StudentModel, error)
	Delete(ctx context.Context, uid uint) (StudentModel, error)
	FindByName(ctx context.Context, studentName string) ([]StudentModel, error)
}
