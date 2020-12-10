package fund

import "time"

// Date is the date of net value of fund
type Date time.Time

// ToDate converts string to Date type
func ToDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return Date(t), nil
}

func (a Date) String() string {
	return time.Time(a).Format("2006-01-02")
}

// IsSame checks if the moment is in date
func (a Date) IsSame(b Moment) bool {
	y1, m1, d1 := time.Time(a).Date()
	y2, m2, d2 := time.Time(b).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsBetween checks if the date is between date range
func (a Date) IsBetween(from, to Date) bool {
	return time.Time(a).After(time.Time(from)) && time.Time(a).Before(time.Time(to))
}
