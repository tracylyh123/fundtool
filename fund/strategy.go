package fund

import (
	"log"
)

// TradeStrategy is a function post by custom
type TradeStrategy func(f *UserFund, total, current Money) error

// Custom is basic info for custom
type Custom struct {
	totalMoney   Money
	currentMoney Money
	funds        []UserFund
}

// Use uses trade strategy of custom
func (c *Custom) Use(s TradeStrategy) {
	for _, f := range c.funds {
		before := f.CalcCurrentMoney()
		err := s(&f, c.totalMoney, c.currentMoney)
		if err != nil {
			log.Printf("cannot apply custom strategy for fund %s, reason: %v", f.Code(), err)
			continue
		}
		c.currentMoney += f.CalcCurrentMoney() - before
		Save(&f)
	}
}

func example(f *UserFund, total, current Money) error {
	e := f.Estval()
	t := f.Trend()
	if e.Status() == Spike {
		rate, err := f.CalcEarningRate()
		if err != nil {
			return err
		}
		if len(e)-1 < 0 || len(t)-1 < 0 || t[len(t)-1].Price == 0 {
			return nil
		}
		delta := e[len(e)-1].IncrRate(t[len(t)-1].Price)
		if rate+delta < 0.1 {
			return nil
		}
		if e.Status() == Spike {
			var amount Share
			if f.TotalMoney() < total/3 {
				amount = f.CurrentShare()
			} else {
				amount = f.CurrentShare() / 2
			}
			if err := f.Sell(amount); err != nil {
				return err
			}
		}
	} else if e.Status() == Drop {
		if len(t) < 3 {
			return nil
		}
		for i := 0; i < 3; i++ {
			if t[len(t)-i].Price >= t[len(t)-i-1].Price {
				return nil
			}
		}
		amount := total / 10
		if amount > current {
			return nil
		}
		if err := f.Buy(amount); err != nil {
			return err
		}
	}
	return nil
}
