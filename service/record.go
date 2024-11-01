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
	if sess != "true" {
		return ErrIDK
	}

	s.rdb.Del(ctx, CreateRecordSessionPrefix(record.UserID))

	today := time.Now().Format("2006-01-02")

	recordCount, err := s.rdb.Get(ctx, RecordsPerDayPrefix(record.UserID, today)).Int()
	if err != nil && err != redis.Nil {
		return err
	}

	if recordCount >= 8 {
		return ErrDailyLimitReached
	}

	if err := s.repo.Record.Create(record); err != nil {
		return err
	}

	if err := s.rdb.Incr(ctx, RecordsPerDayPrefix(record.UserID, today)).Err(); err != nil {
		return err
	}

	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return s.rdb.ExpireAt(ctx, RecordsPerDayPrefix(record.UserID, today), midnight).Err()
}

func (s *RecordService) FindByID(id int64) (*model.Record, error) {
	record, err := s.repo.Record.FindByID(id)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *RecordService) FindWithinMonth(userID int64) ([]*model.Record, error) {
	recordsCache, err := s.rdb.Get(ctx, RecordsWithinMonthPrefix(userID)).Result()
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

	recordsDB, err := s.repo.Record.FindWithinMonth(userID)
	if err != nil {
		return nil, err
	}

	recordsJSON, err := json.Marshal(recordsDB)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, RecordsWithinMonthPrefix(userID), recordsJSON, time.Minute).Err(); err != nil {
		return nil, err
	}

	return recordsDB, nil
}

func (s *RecordService) FindWithinWeek(userID int64) ([]*model.Record, error) {
	recordsCache, err := s.rdb.Get(ctx, RecordsWithinWeekPrefix(userID)).Result()
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

	recordsDB, err := s.repo.Record.FindWithinWeek(userID)
	if err != nil {
		return nil, err
	}

	recordsJSON, err := json.Marshal(recordsDB)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, RecordsWithinWeekPrefix(userID), recordsJSON, time.Minute).Err(); err != nil {
		return nil, err
	}

	return recordsDB, nil
}

func (s *RecordService) FindWithinDay(userID int64) ([]*model.Record, error) {
	records, err := s.repo.Record.FindWithinDay(userID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *RecordService) FindLast(userID int64) (*model.Record, error) {
	record, err := s.repo.Record.FindLast(userID)
	if err != nil {
		return nil, err
	}
	
	return record, nil
}

func (s *RecordService) Search(userID int64, query string) ([]*model.Record, error) {
	records, err := s.repo.Record.Search(userID, query)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *RecordService) Delete(id int64) error {
	return s.repo.Record.Delete(id)
}
