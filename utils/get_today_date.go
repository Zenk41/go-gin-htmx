package utils

import(
	"time"
)

func GetTodayDate()string{
	now := time.Now()

    // Format the date as YYYY-MM-DD
    formattedDate := now.Format("2006-01-02")
	return formattedDate
}