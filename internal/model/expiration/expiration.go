package expiration

import (
	"time"
)

type Expiration int

const (
	HOUR Expiration = iota
	DAY
	WEEK
	HALFMONTH
	MONTH
	HALFYEAR
)

func expirationTypesNameToValueMap() map[string]Expiration {
	return map[string]Expiration{
		"HOUR":      HOUR,
		"DAY":       DAY,
		"WEEK":      WEEK,
		"HALFMONTH": HALFMONTH,
		"MONTH":     MONTH,
		"HALFYEAR":  HALFYEAR,
	}
}

func GetExpiration(exp string, now time.Time) time.Time {
	d := expirationTypesNameToValueMap()
	var expiredAt time.Time
	switch d[exp] {
	case HOUR:
		expiredAt = now.Add(time.Hour)
	case DAY:
		expiredAt = now.Add(24 * time.Hour)
	case WEEK:
		expiredAt = now.Add(7 * 24 * time.Hour)
	case HALFMONTH:
		expiredAt = now.Add(15 * 24 * time.Hour)
	case MONTH:
		expiredAt = now.Add(30 * 24 * time.Hour)
	case HALFYEAR:
		expiredAt = now.Add(6 * 30 * 24 * time.Hour)
	default:
		expiredAt = now.Add(24 * time.Hour)
	}
	return expiredAt
}
