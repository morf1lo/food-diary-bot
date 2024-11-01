package service

import "fmt"

const (
	createRecordSessionPrefix = "create-record-session:%d" // user ID
	recordsWithinMonthPrefix = "records-within-month:%d" // user ID
	recordsWithinWeekPrefix = "records-within-week:%d" // user ID
	recordsPerDayPrefix = "records-per-day:%d:%s" // user ID, date
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

func RecordsPerDayPrefix(userID int64, date string) string {
	return fmt.Sprintf(recordsPerDayPrefix, userID, date)
}
