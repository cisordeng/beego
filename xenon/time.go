package xenon

import "time"

type Time struct {
	time.Time
}

func NewTime(t ...time.Time) Time {
	if len(t) > 0 {
		return Time{t[0]}
	}
	return Time{time.Now()}
}

func (this Time) Format(format ...string) string {
	if this.IsZero() {
		return ""
	}
	if len(format) > 0 {
		return this.Time.Format(format[0])
	}
	return this.Time.Format("2006-01-02 15:04:05")
}