package fund

import (
	"encoding/json"
	"fmt"
)

var events = []string{"Bought", "Sold"}

const (
	boughtIndex = 0
	soldIndex   = 1
)

// Event is the trade event which can be tracked
type Event interface {
	Name() string
}

type eventCommon struct {
	Netval Netval
	Moment Moment
}

type jsonCommon struct {
	Moment string `json:"moment"`
	Netval struct {
		Price string `json:"price"`
		Date  string `json:"date"`
	} `json:"netval"`
}

func (a eventCommon) mapTo(b *jsonCommon) {
	b.Moment = fmt.Sprintf("%s", a.Moment)
	b.Netval.Price = fmt.Sprintf("%s", a.Netval.Price)
	b.Netval.Date = fmt.Sprintf("%s", a.Netval.Date)
	return
}

func (a *eventCommon) mapFrom(b jsonCommon) (err error) {
	a.Moment, err = ToMoment(b.Moment)
	if err != nil {
		return
	}
	a.Netval.Price, err = ToPrice(b.Netval.Price)
	if err != nil {
		return
	}
	a.Netval.Date, err = ToDate(b.Netval.Date)
	if err != nil {
		return
	}
	return
}

// Bought is the event which represents bought into a fund
type Bought struct {
	Money Money
	eventCommon
}

// MarshalJSON converts Bought event to json
func (b *Bought) MarshalJSON() ([]byte, error) {
	var c jsonCommon
	b.mapTo(&c)
	return json.Marshal(&struct {
		Money string `json:"money"`
		jsonCommon
	}{
		fmt.Sprintf("%s", b.Money),
		c,
	})
}

// UnmarshalJSON converts json to Bought event
func (b *Bought) UnmarshalJSON(data []byte) error {
	event := &struct {
		Money string `json:"money"`
		jsonCommon
	}{}
	var err error

	if err = json.Unmarshal(data, event); err != nil {
		return err
	}
	if b.Money, err = ToMoney(event.Money); err != nil {
		return err
	}
	if err = b.mapFrom(event.jsonCommon); err != nil {
		return err
	}
	return nil
}

// CalcShare returns the share which has been bought into
func (b *Bought) CalcShare() Share {
	return Share(float64(b.Money) / float64(b.Netval.Price) * 10000)
}

// Name returns string of event
func (Bought) Name() string {
	return events[boughtIndex]
}

// Sold is the event which represents sold out from a fund
type Sold struct {
	Share Share
	eventCommon
}

// MarshalJSON converts Sold event to json
func (s *Sold) MarshalJSON() ([]byte, error) {
	var c jsonCommon
	s.mapTo(&c)
	return json.Marshal(&struct {
		Share string `json:"share"`
		jsonCommon
	}{
		fmt.Sprintf("%s", s.Share),
		c,
	})
}

// UnmarshalJSON converts json to Sold event
func (s *Sold) UnmarshalJSON(data []byte) error {
	event := &struct {
		Share string `json:"share"`
		jsonCommon
	}{}
	var err error

	if err = json.Unmarshal(data, event); err != nil {
		return err
	}
	s.Share, err = ToShare(event.Share)
	s.mapFrom(event.jsonCommon)
	if err != nil {
		return err
	}
	return nil
}

// CalcMoney returns the money which has been sold out
func (s *Sold) CalcMoney() Money {
	return Money(int(s.Share) * int(s.Netval.Price) / 10000)
}

// Name returns string of event
func (Sold) Name() string {
	return events[soldIndex]
}
