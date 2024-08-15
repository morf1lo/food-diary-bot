package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morf1lo/food-diary-bot/model"
	"github.com/morf1lo/food-diary-bot/repository"
	"github.com/redis/go-redis/v9"
)

type RecordService struct {
	repo *repository.Repository
	rdb *redis.Client
}

var ctx = context.Background()

func NewRecordService(repo *repository.Repository, rdb *redis.Client) *RecordService {
	return &RecordService{
		repo: repo,
		rdb: rdb,
	}
}

func (s *RecordService) RequestToAdd(userID int64) error {
	return s.rdb.Set(context.TODO(), CreateRecordSessionPrefix(userID), "true", time.Minute * 5).Err()
}

func (s *RecordService) Create(record *model.Record) error {
	sess := s.rdb.Get(ctx, CreateRecordSessionPrefix(record.UserID)).Val()
	if sess == "true" {
		s.rdb.Del(ctx, CreateRecordSessionPrefix(record.UserID))
		return s.repo.Record.Create(record)
	}

	return errIDK
}

func (s *RecordService) FindByID(id int64) (*model.Record, error) {
	record, err := s.repo.Record.FindByID(id)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *RecordService) FindWithinMonth(telegramID int64) ([]*model.Record, error) {
	recordsCache, err := s.rdb.Get(ctx, RecordsWithinMonthPrefix(telegramID)).Result()
	if err == nil {
		var records []*model.Record
		if err := json.Unmarshal([]byte(recordsCache), &records); err != nil {
			return nil, err
		}
		return records, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	recordsDB, err := s.repo.Record.FindWithinMonth(telegramID)
	if err != nil {
		return nil, err
	}

	recordsJSON, err := json.Marshal(recordsDB)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, RecordsWithinMonthPrefix(telegramID), recordsJSON, time.Minute).Err(); err != nil {
		return nil, err
	}

	return recordsDB, nil
}

func (s *RecordService) FindWithinWeek(telegramID int64) ([]*model.Record, error) {
	recordsCache, err := s.rdb.Get(ctx, RecordsWithinWeekPrefix(telegramID)).Result()
	if err == nil {
		var records []*model.Record
		if err := json.Unmarshal([]byte(recordsCache), &records); err != nil {
			return nil, err
		}
		return records, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	recordsDB, err := s.repo.Record.FindWithinWeek(telegramID)
	if err != nil {
		return nil, err
	}

	recordsJSON, err := json.Marshal(recordsDB)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, RecordsWithinWeekPrefix(telegramID), recordsJSON, time.Minute).Err(); err != nil {
		return nil, err
	}

	return recordsDB, nil
}

func (s *RecordService) FindWithinDay(telegramID int64) ([]*model.Record, error) {
	records, err := s.repo.Record.FindWithinDay(telegramID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *RecordService) FindLast(telegramID int64) (*model.Record, error) {
	record, err := s.repo.Record.FindLast(telegramID)
	if err != nil {
		return nil, err
	}
	
	return record, nil
}

func (s *RecordService) Search(telegramID int64, query string) ([]*model.Record, error) {
	records, err := s.repo.Record.Search(telegramID, query)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *RecordService) Delete(id int64) error {
	return s.repo.Record.Delete(id)
}
