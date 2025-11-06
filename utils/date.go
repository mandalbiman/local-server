package utils

import "time"

func FormatDate(t *time.Time, layout string) string {
	return t.Format(layout)
}

func FormatDateTime(t *time.Time) string {
	return FormatDate(t, "02-01-2006 03:04 PM")
}

func FormatDateAsStr(t *time.Time) string {
	return FormatDate(t, "02-01-2006")
}
