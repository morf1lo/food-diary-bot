package service

import (
	"context"
	"strconv"
	"time"

	"github.com/morf1lo/food-diary-bot/model"
	"github.com/morf1lo/food-diary-bot/repository"
	"github.com/redis/go-redis/v9"
)

const createRecordSessionPrefix = "create-record-session-"

type RecordService struct {
	repo *repository.Repository
	rdb *redis.Client
}

func NewRecordService(repo *repository.Repository, rdb *redis.Client) *RecordService {
	return &RecordService{
		repo: repo,
		rdb: rdb,
	}
}

func (s *RecordService) RequestToAdd(userID string) error {
	return s.rdb.Set(context.TODO(), createRecordSessionPrefix + userID, "true", time.Minute * 5).Err()
}

func (s *RecordService) Create(record *model.Record) error {
	userIDString := strconv.Itoa(int(record.UserID))
	sess := s.rdb.Get(context.TODO(), createRecordSessionPrefix + userIDString).Val()
	if sess == "true" {
		s.rdb.Del(context.TODO(), createRecordSessionPrefix + userIDString)

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

func (s *RecordService) FindAll(telegramID int64) ([]*model.Record, error) {
	records, err := s.repo.Record.FindAll(telegramID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *RecordService) FindByMonth(telegramID int64) ([]*model.Record, error) {
	records, err := s.repo.Record.FindByMonth(telegramID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *RecordService) FindByWeek(telegramID int64) ([]*model.Record, error) {
	records, err := s.repo.Record.FindByWeek(telegramID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *RecordService) FindByDay(telegramID int64) ([]*model.Record, error) {
	records, err := s.repo.Record.FindByDay(telegramID)
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
