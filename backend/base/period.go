package base

import (
	"encoding"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type periodImpl struct {
	Year   int
	Month  int
	Week   int
	Day    int
	Hour   int
	Minute int
	Second int
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *periodImpl) UnmarshalText(text []byte) (err error) {
	*p, err = ParsePeriod(string(text))
	return
}

// MarshalText implements encoding.TextMarshaler.
func (p periodImpl) MarshalText() (text []byte, err error) {
	sb := strings.Builder{}
	sb.WriteByte('P')
	if p.Year != 0 {
		sb.WriteString(strconv.Itoa(p.Year))
		sb.WriteByte('Y')
	}
	if p.Month != 0 {
		sb.WriteString(strconv.Itoa(p.Month))
		sb.WriteByte('M')
	}
	if p.Week != 0 {
		sb.WriteString(strconv.Itoa(p.Week))
		sb.WriteByte('W')
	}
	if p.Day != 0 {
		sb.WriteString(strconv.Itoa(p.Day))
		sb.WriteByte('D')
	}
	if p.Hour != 0 || p.Minute != 0 || p.Second != 0 {
		sb.WriteByte('T')
		if p.Hour != 0 {
			sb.WriteString(strconv.Itoa(p.Hour))
			sb.WriteByte('H')
		}
		if p.Minute != 0 {
			sb.WriteString(strconv.Itoa(p.Minute))
			sb.WriteByte('M')
		}
		if p.Second != 0 {
			sb.WriteString(strconv.Itoa(p.Second))
			sb.WriteByte('S')
		}
	}
	return []byte(sb.String()), nil
}

var ErrInvalidPeriod = errors.New("not a valid period")

func periodRegexCaptureGroup(suffix string, name string) string {
	return fmt.Sprintf("(?:(?P<%s>[0-9]+)%s)?", name, suffix)
}

var periodRegex = regexp.MustCompile(fmt.Sprintf(`^P%s%s%s%s(?:T%s%s%s)?$`,
	periodRegexCaptureGroup("Y", "year"),
	periodRegexCaptureGroup("M", "month"),
	periodRegexCaptureGroup("W", "week"),
	periodRegexCaptureGroup("D", "day"),
	periodRegexCaptureGroup("H", "hour"),
	periodRegexCaptureGroup("M", "minute"),
	periodRegexCaptureGroup("S", "second"))) // D, T, H, M, S$`)

func ParsePeriod(s string) (p periodImpl, err error) {
	match := periodRegex.FindStringSubmatch(s)
	if match == nil {
		return p, ErrInvalidPeriod
	}
	for i, name := range periodRegex.SubexpNames() {
		if i != 0 && name != "" && match[i] != "" {
			switch name {
			case "year":
				p.Year, err = strconv.Atoi(match[i])
			case "month":
				p.Month, err = strconv.Atoi(match[i])
			case "week":
				p.Week, err = strconv.Atoi(match[i])
			case "day":
				p.Day, err = strconv.Atoi(match[i])
			case "hour":
				p.Hour, err = strconv.Atoi(match[i])
			case "minute":
				p.Minute, err = strconv.Atoi(match[i])
			case "second":
				p.Second, err = strconv.Atoi(match[i])
			}
			if err != nil {
				return p, fmt.Errorf("%w: invalid %s: %w", ErrInvalidPeriod, name, err)
			}
		}
	}
	return
}

var _ encoding.TextMarshaler = periodImpl{}
var _ encoding.TextUnmarshaler = (*periodImpl)(nil)
