package decimal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidPeriod = errors.New("invalid billing period format")

// Period represents a recurring billing cycle (e.g. 1W, 1M, 3M, 6M, 1Y, 2Y)
type Period struct {
	Unit  string // "W", "M", "Y"
	Qty   int    // 1, 3, 6, 12, etc.
	Title string
}

func ParsePeriod(code string) (*Period, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" || code == "1" || code == "ONETIME" {
		return &Period{Unit: "O", Qty: 1, Title: "One Time"}, nil
	}
	if code == "FREE" {
		return &Period{Unit: "F", Qty: 0, Title: "Free"}, nil
	}

	if len(code) < 2 {
		return nil, ErrInvalidPeriod
	}

	unit := string(code[len(code)-1])
	qtyStr := code[:len(code)-1]

	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty <= 0 {
		return nil, ErrInvalidPeriod
	}

	var title string
	switch unit {
	case "W":
		if qty == 1 {
			title = "Weekly"
		} else {
			title = fmt.Sprintf("Every %d Weeks", qty)
		}
	case "M":
		switch qty {
		case 1:
			title = "Monthly"
		case 3:
			title = "Quarterly"
		case 6:
			title = "Semi-Annually"
		default:
			title = fmt.Sprintf("Every %d Months", qty)
		}
	case "Y":
		switch qty {
		case 1:
			title = "Annually"
		case 2:
			title = "Biennially"
		case 3:
			title = "Triennially"
		default:
			title = fmt.Sprintf("Every %d Years", qty)
		}
	default:
		return nil, ErrInvalidPeriod
	}

	return &Period{
		Unit:  unit,
		Qty:   qty,
		Title: title,
	}, nil
}

func (p *Period) CalculateNextDueDate(fromDate time.Time) time.Time {
	switch p.Unit {
	case "W":
		return fromDate.AddDate(0, 0, p.Qty*7)
	case "M":
		return fromDate.AddDate(0, p.Qty, 0)
	case "Y":
		return fromDate.AddDate(p.Qty, 0, 0)
	default:
		return fromDate
	}
}
