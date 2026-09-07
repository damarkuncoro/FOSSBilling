package period

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	UnitDay   = "D"
	UnitWeek  = "W"
	UnitMonth = "M"
	UnitYear  = "Y"

	PeriodWeek         = "1W"
	PeriodMonth        = "1M"
	PeriodQuarter      = "3M"
	PeriodBiannual     = "6M"
	PeriodAnnual       = "1Y"
	PeriodBiennial     = "2Y"
	PeriodTriennial    = "3Y"
	PeriodQuadrennial  = "4Y"
	PeriodQuinquennial = "5Y"
)

var periodRegex = regexp.MustCompile(`^(\d+)([A-Za-z])$`)

type Period struct {
	Code string `json:"code"`
	Qty  int    `json:"qty"`
	Unit string `json:"unit"`
}

func Parse(code string) (*Period, error) {
	code = strings.TrimSpace(code)
	matches := periodRegex.FindStringSubmatch(code)
	if len(matches) != 3 {
		return nil, errors.New("invalid period code format (expected e.g. 1M, 1Y, 3M, 1W)")
	}

	qty, err := strconv.Atoi(matches[1])
	if err != nil || qty <= 0 {
		return nil, errors.New("period quantity must be positive")
	}

	unit := strings.ToUpper(matches[2])
	switch unit {
	case UnitDay:
		if qty > 90 {
			return nil, errors.New("day quantity cannot exceed 90")
		}
	case UnitWeek:
		if qty > 52 {
			return nil, errors.New("week quantity cannot exceed 52")
		}
	case UnitMonth:
		if qty > 24 {
			return nil, errors.New("month quantity cannot exceed 24")
		}
	case UnitYear:
		if qty > 5 {
			return nil, errors.New("year quantity cannot exceed 5")
		}
	default:
		return nil, fmt.Errorf("unknown period unit: %s", unit)
	}

	return &Period{
		Code: fmt.Sprintf("%d%s", qty, unit),
		Qty:  qty,
		Unit: unit,
	}, nil
}

func (p *Period) CalculateExpiration(from time.Time) time.Time {
	switch p.Unit {
	case UnitDay:
		return from.AddDate(0, 0, p.Qty)
	case UnitWeek:
		return from.AddDate(0, 0, p.Qty*7)
	case UnitMonth:
		return from.AddDate(0, p.Qty, 0)
	case UnitYear:
		return from.AddDate(p.Qty, 0, 0)
	default:
		return from.AddDate(0, p.Qty, 0)
	}
}

func (p *Period) GetTitle() string {
	switch p.Code {
	case PeriodWeek:
		return "Every Week"
	case PeriodMonth:
		return "Every Month"
	case PeriodQuarter:
		return "Every 3 Months"
	case PeriodBiannual:
		return "Every 6 Months"
	case PeriodAnnual:
		return "Every Year"
	case PeriodBiennial:
		return "Every 2 Years"
	case PeriodTriennial:
		return "Every 3 Years"
	default:
		return fmt.Sprintf("Every %d %s", p.Qty, p.Unit)
	}
}
