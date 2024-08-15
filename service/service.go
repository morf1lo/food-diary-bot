package service

import (
	"github.com/morf1lo/food-diary-bot/model"
	"github.com/morf1lo/food-diary-bot/repository"
	"github.com/redis/go-redis/v9"
)

type Record interface {
	RequestToAdd(userID int64) error
	Create(record *model.Record) error
	FindByID(id int64) (*model.Record, error)
	FindWithinMonth(telegramID int64) ([]*model.Record, error)
	FindWithinWeek(telegramID int64) ([]*model.Record, error)
	FindWithinDay(telegramID int64) ([]*model.Record, error)
	FindLast(telegramID int64) (*model.Record, error)
	Search(telegramID int64, query string) ([]*model.Record, error)
	Delete(id int64) error
}

type Service struct {
	Record
}

func New(repo *repository.Repository, rdb *redis.Client) *Service {
	return &Service{
		Record: NewRecordService(repo, rdb),
	}
}
