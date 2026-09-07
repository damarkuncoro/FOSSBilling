package decimal

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Money represents monetary amounts in hundredths of a cent (4 decimal places precision)
type Money int64

func FromFloat(f float64) Money {
	return Money(math.Round(f * 10000))
}

func FromString(s string) (Money, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return FromFloat(f), nil
}

func (m Money) ToFloat() float64 {
	return float64(m) / 10000.0
}

func (m Money) String() string {
	val := int64(m)
	sign := ""
	if val < 0 {
		sign = "-"
		val = -val
	}
	// Round to nearest 2 decimal places (hundredths of a dollar = cents)
	cents := (val + 50) / 100
	dollars := cents / 100
	remainderCents := cents % 100
	return fmt.Sprintf("%s%d.%02d", sign, dollars, remainderCents)
}

func (m Money) FormatPrecision(decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, m.ToFloat())
}

func (m *Money) UnmarshalJSON(data []byte) error {
	// 1. Try as string (e.g. "9.99")
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		val, err := FromString(s)
		if err != nil {
			return err
		}
		*m = val
		return nil
	}

	// 2. Try as float (e.g. 9.99 or 100)
	// We treat ALL numbers in JSON as float (dollars) for consistency
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*m = FromFloat(f)
		return nil
	}

	return fmt.Errorf("invalid money value: %s", string(data))
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToFloat())
}
