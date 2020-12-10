package fetcher

import (
	"encoding/json"
	"fmt"
	"strings"
)

// QueryTranslator transforms data comfortable with database
type QueryTranslator interface {
	startParser(out chan<- interface{}, in <-chan *body)
	startBuilder(out chan<- *query, in <-chan interface{})
}

type fundEst struct {
	Code      string `json:"fundcode"`
	EstVal    string `json:"gsz"`
	EstChange string `json:"gszzl"`
	EstTime   string `json:"gztime"`
}

type fundNet struct {
	NetDate string `json:"x"`
	NetVal  string `json:"y"`
}

// FundEstQueryTranslator transforms estimated data comfortable with database
type FundEstQueryTranslator struct{}

func (f *FundEstQueryTranslator) startParser(out chan<- interface{}, in <-chan *body) {
	parser(out, in, func(b *body) (interface{}, error) {
		var f fundEst
		s := fmt.Sprintf("%s", b.payload)
		err := json.Unmarshal([]byte(s[8:len(s)-2]), &f)
		if err != nil {
			return nil, fmt.Errorf("cannot parse body from %s, reason: %v", b.url, err)
		}
		return &f, nil
	})
}

func (f *FundEstQueryTranslator) startBuilder(out chan<- *query, in <-chan interface{}) {
	var buf strings.Builder
	batchNum := 1000
	vals := make([]interface{}, 0, batchNum*4)
	clause := "INSERT INTO `fundest` (`code`, `estval`, `estchange`, `esttime`) VALUES"
	n := 0
	defer close(out)
	for i := range in {
		if f, ok := i.(*fundEst); ok {
			if n == 0 {
				buf.WriteString(clause)
			} else {
				buf.WriteString(",")
			}
			buf.WriteString("(?, ?, ?, ?)")
			vals = append(vals, f.Code, f.EstVal, f.EstChange, f.EstTime)
			n++
			if n == batchNum {
				q := query{sql: buf.String(), params: make([]interface{}, len(vals))}
				copy(q.params, vals)
				out <- &q
				buf.Reset()
				vals = vals[:0]
				n = 0
			}
		} else {
			panic(fmt.Sprintf("unexpected type: %v", i))
		}
	}
	if len(vals) > 0 {
		out <- &query{sql: buf.String(), params: vals}
	}
}

type fundNetHistory []fundNet

// FundNetHistoryQueryTranslator transforms net value data comfortable with database
type FundNetHistoryQueryTranslator struct{}

func (f *FundNetHistoryQueryTranslator) startParser(out chan<- interface{}, in <-chan *body) {
	parser(out, in, func(b *body) (interface{}, error) {
		var f fundNetHistory
		s := fmt.Sprintf("%s", b.payload)
		begin := strings.Index(s, "Data_netWorthTrend")
		t := s[begin:]
		end := strings.Index(t, ";")
		err := json.Unmarshal([]byte(t[:end]), &f)
		if err != nil {
			return nil, fmt.Errorf("cannot parse body from %s, reason: %v", b.url, err)
		}
		return &f, nil
	})
}

func (f *FundNetHistoryQueryTranslator) startBuilder(out chan<- *query, in <-chan interface{}) {
	clause := "INSERT IGNORE INTO `fundnet` (`netval`, `netdate`) VALUES"
	defer close(out)
	for i := range in {
		if h, ok := i.(*fundNetHistory); ok {
			var buf strings.Builder
			var vals []interface{}
			buf.WriteString(clause)
			for n, f := range *h {
				if n != 0 {
					buf.WriteString(",")
				}
				buf.WriteString("(?, ?)")
				vals = append(vals, f.NetVal, f.NetDate)
			}
			out <- &query{sql: buf.String(), params: vals}
		} else {
			panic(fmt.Sprintf("unexpected type: %v", i))
		}
	}
}
