package service

import "fmt"

const (
	createRecordSessionPrefix = "create-record-session:%d" // user ID
	recordsWithinMonthPrefix = "records-within-month:%d" // user ID
	recordsWithinWeekPrefix = "records-within-week:%d" // user ID
)

func CreateRecordSessionPrefix(userID int64) string {
	return fmt.Sprintf(createRecordSessionPrefix, userID)
}

func RecordsWithinMonthPrefix(userID int64) string {
	return fmt.Sprintf(recordsWithinMonthPrefix, userID)
}

func RecordsWithinWeekPrefix(userID int64) string {
	return fmt.Sprintf(recordsWithinWeekPrefix, userID)
}
