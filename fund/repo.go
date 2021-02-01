package fund

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/tracylyh123/fundtool/helper"
)

var conf helper.Config = helper.Global.Config

// Save persists all evnets on MyFund
func Save(f *UserFund) error {
	var buf strings.Builder
	vals := make([]interface{}, 0)
	prefix := "INSERT INTO `fundevents` (`event_type`, `code`, `payload`, `version`, `created_by`, `created_at`) VALUES"
	for n, event := range f.events {
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		if n == 0 {
			buf.WriteString(prefix)
		} else {
			buf.WriteString(",")
		}
		buf.WriteString("(?, ?, ?, ?, ?, NOW())")
		vals = append(vals, event.Name(), f.code, string(data), f.version+n, f.userID)
	}
	db, err := sql.Open(conf.DB.Driver, conf.DB.DNS)
	if err != nil {
		return err
	}
	defer db.Close()
	if buf.Len() != 0 {
		_, err = db.Exec(buf.String(), vals...)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReplayEvents replays all events on user's fund
func ReplayEvents(f *UserFund) error {
	db, err := sql.Open(conf.DB.Driver, conf.DB.DNS)
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("SELECT event_type, payload FROM `fundevents` WHERE `created_by`=? AND `code`=?", f.userID, f.code)
	if err != nil {
		return err
	}
	for rows.Next() {
		var item struct {
			name    string
			payload []byte
		}
		err := rows.Scan(&item.name, &item.payload)
		if err != nil {
			return err
		}
		var ev Event
		switch item.name {
		case events[boughtIndex]:
			ev = &Bought{}
		case events[soldIndex]:
			ev = &Sold{}
		default:
			panic(fmt.Sprintf("unexpected event name: %s", item.name))
		}
		err = json.Unmarshal(item.payload, ev)
		if err != nil {
			return err
		}
		f.Apply(ev, false)
	}
	return nil
}

// Find finds UserFund in DB
func Find(userID int, code string) (*UserFund, error) {
	db, err := sql.Open(conf.DB.Driver, conf.DB.DNS)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT `netval`, `netdate` FROM `fundnet` WHERE `code`=? ORDER BY `netdate` DESC LIMIT 1", code)
	var item struct {
		netval  string
		netdate string
	}
	err = row.Scan(&item.netval, &item.netdate)
	if err != nil {
		return nil, err
	}
	v, err := ToPrice(item.netval)
	if err != nil {
		return nil, err
	}
	t, err := ToDate(item.netdate)
	if err != nil {
		return nil, err
	}
	f := Fund{code: code, netval: Netval{Price: v, Date: t}}
	return &UserFund{Fund: f, userID: userID}, nil
}

// LoadTrend loads trend for fund
func LoadTrend(fund *UserFund, from, to string) error {
	trend, err := FindTrend(fund.code, from, to)
	if err != nil {
		return err
	}
	fund.trend = trend
	return nil
}

// FindTrend finds trend for fund
func FindTrend(code, from, to string) (Trend, error) {
	db, err := sql.Open(conf.DB.Driver, conf.DB.DNS)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT `netval`, `netdate` FROM `fundnet` WHERE `code`=? AND `netdate` BETWEEN ? AND ? ORDER BY `netdate` DESC", code, from, to)
	if err != nil {
		return nil, err
	}
	var trend Trend
	for rows.Next() {
		var item struct {
			netval  string
			netdate string
		}
		err = rows.Scan(&item.netval, &item.netdate)
		if err != nil {
			return nil, err
		}
		v, err := ToPrice(item.netval)
		if err != nil {
			return nil, err
		}
		t, err := ToDate(item.netdate)
		if err != nil {
			return nil, err
		}
		trend = append(trend, Netval{Price: v, Date: t})
	}
	return trend, nil
}
