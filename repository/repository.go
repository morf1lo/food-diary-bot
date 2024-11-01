package repository

import (
	"github.com/morf1lo/food-diary-bot/model"
	"gorm.io/gorm"
)

type Record interface {
	Create(record *model.Record) error
	FindByID(id int64) (*model.Record, error)
	FindWithinMonth(userID int64) ([]*model.Record, error)
	FindWithinWeek(userID int64) ([]*model.Record, error)
	FindWithinDay(userID int64) ([]*model.Record, error)
	FindLast(userID int64) (*model.Record, error)
	Search(userID int64, query string) ([]*model.Record, error)
	Delete(id int64) error
}

type Repository struct {
	Record
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		Record: NewRecordRepo(db),
	}
}
