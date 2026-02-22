package timeutil

import "time"

func NowAsiaJakarta() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}

func FormatISO8601(t time.Time) string {
	return t.Format(time.RFC3339)
}
