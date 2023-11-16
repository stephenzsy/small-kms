package base

import (
	"time"

	"github.com/stephenzsy/small-kms/backend/utils/caldur"
)

type periodImpl = caldur.CalendarDuration

// Deprecated: use caldur.CalendarDuration instead.
func AddPeriod(t time.Time, p Period) time.Time {
	return t.UTC().AddDate(p.Year, p.Month, p.Day).Add(
		time.Duration(p.Hour)*time.Hour +
			time.Duration(p.Minute)*time.Minute +
			time.Duration(p.Second)*time.Second)
}
