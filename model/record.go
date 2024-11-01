package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Record struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"column:user_id;unique;not null" validate:"required"`
	Body      string    `gorm:"column:body;not null" validate:"required,min=1,max=128"`
	DateAdded time.Time `gorm:"column:date_added;type:timestamp(0) without time zone;not null;autoCreateTime"`
}

func (r *Record) Validate() error {
	return validator.New().Struct(r)
}
