package service

import "errors"

var (
	ErrIDK = errors.New("idk")
	ErrDailyLimitReached = errors.New("you have reached the daily limit of 8 records")
)
