package decimal

import (
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
