package fund

import (
	"fmt"
	"strconv"
)

// Share = Money / Price
type Share int

func (s Share) String() string {
	return fmt.Sprintf("%.2f", float64(s)/100)
}

// ToShare converts string to Share type
func ToShare(s string) (Share, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert string: %s to share", s)
	}
	return Share((f * 100) + 0.5), nil
}
