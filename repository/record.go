package repository

import (
	"github.com/morf1lo/food-diary-bot/model"
	"gorm.io/gorm"
)

type RecordRepo struct {
	db *gorm.DB
}

func NewRecordRepo(db *gorm.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

func (r *RecordRepo) Create(record *model.Record) error {
	return r.db.Save(record).Error
}

func (r *RecordRepo) FindByID(id int64) (*model.Record, error) {
	var record model.Record
	if err := r.db.Where("id = ?", id).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *RecordRepo) FindWithinMonth(telegramID int64) ([]*model.Record, error) {
	var records []*model.Record
	if err := r.db.Raw("SELECT * FROM records r WHERE r.user_id = $1 AND r.date_added >= CURRENT_TIMESTAMP - INTERVAL '1 month'", telegramID).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *RecordRepo) FindWithinWeek(telegramID int64) ([]*model.Record, error) {
	var records []*model.Record
	if err := r.db.Raw("SELECT * FROM records r WHERE r.user_id = $1 AND r.date_added >= CURRENT_TIMESTAMP - INTERVAL '1 week'", telegramID).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *RecordRepo) FindWithinDay(telegramID int64) ([]*model.Record, error) {
	var records []*model.Record
	if err := r.db.Raw("SELECT * FROM records r WHERE r.user_id = $1 AND r.date_added >= CURRENT_DATE", telegramID).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *RecordRepo) FindLast(telegramID int64) (*model.Record, error) {
	var record model.Record
	if err := r.db.Raw("SELECT * FROM records r WHERE r.user_id = $1 ORDER BY r.date_added DESC LIMIT 1", telegramID).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *RecordRepo) Search(telegramID int64, query string) ([]*model.Record, error) {
	var records []*model.Record
	if err := r.db.Raw("SELECT * FROM records r WHERE r.user_id = $1 AND r.date_added >= CURRENT_TIMESTAMP - INTERVAL '1 month' AND r.body ILIKE $2", telegramID, "%"+query+"%").Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *RecordRepo) Delete(id int64) error {
	return r.db.Delete(&model.Record{}, id).Error
}
