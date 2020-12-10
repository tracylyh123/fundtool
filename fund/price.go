package fund

import (
	"fmt"
	"strconv"
)

// Price is unit price of fund
type Price int

// IncrRate returns the rate of increase
func (p Price) IncrRate(q Price) float64 {
	return float64((p - q) / q)
}

func (p Price) String() string {
	return fmt.Sprintf("%.4f", float64(p)/10000)
}

// ToPrice converts string to Price type
func ToPrice(s string) (Price, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert string: %s to price", s)
	}
	return Price((f * 10000) + 0.5), nil
}
