package fund

import "fmt"

// Trend is a set of Netval
type Trend []Netval

// Between returns net value of trend between from and to
func (t Trend) Between(from, to Date) Trend {
	var r Trend
	for _, netval := range t {
		if netval.Date.IsBetween(from, to) {
			r = append(r, netval)
		}
	}
	return r
}

// MA is moving average metric
func (t Trend) MA(days int) (Trend, error) {
	if days <= 0 {
		return nil, fmt.Errorf("days should be a positive integer")
	}
	var r Trend
	for i := len(t) - 1; i+1 >= days; i-- {
		var sum Price = t[i].Price
		for j := i - 1; j > i-days; j-- {
			sum += t[j].Price
		}
		avg := int(sum) / days
		r = append(r, Netval{
			Price: Price(avg),
			Date:  t[i].Date,
		})
	}
	return r, nil
}

// MDD is maximum draw down metric
func (t Trend) MDD() (float64, error) {
	if len(t) == 0 {
		return 0, fmt.Errorf("calculating failed: trend was empty")
	}
	max := t[0].Price
	min := t[0].Price
	for _, n := range t {
		if n.Price > max {
			max = n.Price
		}
		if n.Price < min {
			min = n.Price
		}
	}
	if max == 0 {
		return 0, fmt.Errorf("calculating failed: unexpected maximum")
	}
	return float64((max - min) / max), nil
}
