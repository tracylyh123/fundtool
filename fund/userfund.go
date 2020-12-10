package fund

import (
	"fmt"
	"time"
)

// UserFund is the info which is related to user
type UserFund struct {
	Fund
	userID       int
	totalMoney   Money
	currentShare Share
	events       []Event
	version      int
}

// TotalMoney returns total money of fund
func (f UserFund) TotalMoney() Money {
	return f.totalMoney
}

// CurrentShare returns current share of fund
func (f UserFund) CurrentShare() Share {
	return f.currentShare
}

func (f *UserFund) buy(m Money, t Moment) error {
	if !f.IsTradable(t) {
		return fmt.Errorf("cannot buy into fund %s on date: %s", f.code, f.netval.Date)
	}
	ev := &Bought{m, eventCommon{f.netval, t}}
	f.raise(ev)
	return nil
}

// Buy will raise a Bought event
func (f *UserFund) Buy(m Money) error {
	return f.buy(m, Moment(time.Now()))
}

func (f *UserFund) sell(s Share, t Moment) error {
	if !f.IsTradable(t) {
		return fmt.Errorf("cannot sell out fund %s on date: %s", f.code, f.netval.Date)
	}
	ev := Sold{s, eventCommon{f.netval, t}}
	if s > f.currentShare {
		return fmt.Errorf("cannot sell out invalid share: %s", s)
	}
	f.raise(&ev)
	return nil
}

// Sell will raise a Sold event
func (f *UserFund) Sell(s Share) error {
	return f.sell(s, Moment(time.Now()))
}

// CalcCurrentMoney returns current money, formula: money = share * price
func (f UserFund) CalcCurrentMoney() Money {
	return Money(int(f.currentShare) * int(f.netval.Price) / 10000)
}

// CalcEarning returns current earning, formula: earning = current - total
func (f UserFund) CalcEarning() Money {
	return f.CalcCurrentMoney() - f.totalMoney
}

// CalcEarningRate returns current earning rate, formula rate = current / total
func (f UserFund) CalcEarningRate() (float64, error) {
	if f.totalMoney == 0 {
		return 0, fmt.Errorf("earning rate cannot divided by 0")
	}
	return float64(f.CalcCurrentMoney()) / float64(f.totalMoney), nil
}

func (f *UserFund) raise(ev Event) {
	f.events = append(f.events, ev)
	f.Apply(ev, true)
}

// Apply applys event on UserFund
func (f *UserFund) Apply(ev Event, isNew bool) {
	switch e := ev.(type) {
	case *Bought:
		f.totalMoney += e.Money
		f.currentShare += e.CalcShare()
	case *Sold:
		f.totalMoney -= e.CalcMoney()
		f.currentShare -= e.Share
	default:
		panic(fmt.Sprintf("invalid event type: %s", e))
	}

	if !isNew {
		f.version++
	}
}

// NewUserFund creates a UserFund struct
func NewUserFund(userID int, code string, price Price, date Date) *UserFund {
	return &UserFund{
		Fund: Fund{
			code,
			Netval{
				price,
				date,
			},
			nil,
			nil,
		},
		userID: userID,
	}
}
