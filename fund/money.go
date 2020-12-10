package fund

import (
	"fmt"
	"strconv"
)

// Money = Share * Price
type Money int

func (m Money) String() string {
	return fmt.Sprintf("%.2f", float64(m)/100)
}

// ToMoney converts string to Money type
func ToMoney(s string) (Money, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert string: %s to money", s)
	}
	return Money((f * 100) + 0.5), nil
}
