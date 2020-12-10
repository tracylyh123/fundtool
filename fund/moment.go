package fund

import "time"

// Moment is the moment which the trade event occurred
type Moment time.Time

// ToMoment converts string to Moment type
func ToMoment(s string) (Moment, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Moment{}, err
	}
	return Moment(t), nil
}

func (m Moment) String() string {
	return time.Time(m).Format(time.RFC3339)
}
