package handlers

import (
	"time"
)

func DateToTime(date string) (time.Time, error) {
	layout := "2006-01-02" // The layout specifies the format of the input date string
	return time.Parse(layout, date)
}
func getDateOnly(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
