package student

import (
	"context"
	"errors"
	"myschool/pkg/model"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// implmentation model of Student
type Student struct {
	UID       uint                  `gorm:"primaryKey;autoIncrement"`
	CreatedAt int                   `gorm:"autoCreateTime"` //store as second unix timestamp to have compatibility with SQLite
	UpdatedAt int                   `gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `gorm:"index"`
	Name      string
	Gender    string
}

// Safe checker to know if this file already implements the model interface correctly or not
var _ (model.StudentModel) = (*Student)(nil)

func (s *Student) GetUID() uint {
	return s.UID
}

func (s *Student) GetCreatedAt() int {
	return s.CreatedAt
}

func (s *Student) GetUpdatedAt() int {
	return s.UpdatedAt
}

func (s *Student) GetDeletedAt() int {
	return int(s.DeletedAt)
}

func (s *Student) GetName() string {
	return s.Name
}

func (s *Student) GetGender() string {
	return s.Gender
}

// unexportable
type repoPrivate struct {
	db *gorm.DB
}

// Safe checker to know if this file already implements the interface correctly or not
//var _ (model.StudentRepo) = (*repoPrivate)(nil)
//Note: just for example. Not needed if we have constructor like below, NewRepo returning interface{}

// public constructor, requires actual database object
func NewRepo(db *gorm.DB) model.StudentRepo {
	return &repoPrivate{db: db}
}

func (repo *repoPrivate) Create(ctx context.Context, args model.StudentCreationParam) (model.StudentModel, error) {

	studentObject := Student{Name: args.Name, Gender: args.Gender}

	result := repo.db.Create(&studentObject)
	return &studentObject, result.Error
}

func (repo *repoPrivate) Read(ctx context.Context, uid uint) (model.StudentModel, error) {
	studentObject := Student{UID: uid}

	result := repo.db.Find(&studentObject)
	if studentObject.CreatedAt < 1 {
		return &studentObject, errors.New("user not found")
	}

	return &studentObject, result.Error
}

func (repo *repoPrivate) Delete(ctx context.Context, uid uint) (model.StudentModel, error) {
	studentObject := Student{UID: uid}

	result := repo.db.Find(&studentObject)
	if result.Error != nil {
		return &studentObject, result.Error
	}
	if studentObject.CreatedAt < 1 {
		return &studentObject, errors.New("user not found")
	}

	result = repo.db.Delete(&studentObject)
	return &studentObject, result.Error
}

func (repo *repoPrivate) FindByName(ctx context.Context, name string) ([]model.StudentModel, error) {
	var foundStudents []*Student
	result := repo.db.Where("name LIKE ?", "%"+name+"%").Find(&foundStudents)

	models := make([]model.StudentModel, len(foundStudents))
	for i, v := range foundStudents {
		models[i] = model.StudentModel(v)
	}
	return models, result.Error
}
