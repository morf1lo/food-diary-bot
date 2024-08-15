package repository

import (
	"github.com/morf1lo/food-diary-bot/model"
	"gorm.io/gorm"
)

type Record interface {
	Create(record *model.Record) error
	FindByID(id int64) (*model.Record, error)
	FindWithinMonth(telegramID int64) ([]*model.Record, error)
	FindWithinWeek(telegramID int64) ([]*model.Record, error)
	FindWithinDay(telegramID int64) ([]*model.Record, error)
	FindLast(telegramID int64) (*model.Record, error)
	Search(telegramID int64, query string) ([]*model.Record, error)
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
