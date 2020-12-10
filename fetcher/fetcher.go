package fetcher

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/tracylyh123/fundtool/helper"
)

var conf helper.Config = helper.Global.Config

const (
	timeoutSec    = 120
	maxCrawlerNum = 20
)

type body struct {
	payload []byte
	url     string
}

type query struct {
	sql    string
	params []interface{}
}

var wg1, wg2 sync.WaitGroup

func crawl(url string) (*body, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get response from %s, reason: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid http status code %d from %s", resp.StatusCode, url)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot get body from %s, reason: %v", url, err)
	}
	return &body{payload: b, url: url}, nil
}

func startCrawler(out chan<- *body, in <-chan string) <-chan time.Duration {
	elapsed := make(chan time.Duration)
	go func() {
		defer func() {
			wg1.Wait()
			close(out)
			close(elapsed)
		}()
		for i := 0; i < maxCrawlerNum; i++ {
			go func() {
				for url := range in {
					start := time.Now()
					bp, err := crawl(url)
					elapsed <- time.Since(start)
					if err != nil {
						log.Print(err)
						wg2.Done()
					} else {
						out <- bp
					}
					wg1.Done()
				}
			}()
		}
	}()
	return elapsed
}

func parser(out chan<- interface{}, in <-chan *body, handle func(*body) (interface{}, error)) {
	defer func() {
		wg2.Wait()
		close(out)
	}()
	for bp := range in {
		go func(bp *body) {
			parsed, err := handle(bp)
			if err != nil {
				log.Print(err)
			} else {
				out <- parsed
			}
			wg2.Done()
		}(bp)
	}
}

func startTranslator(out chan<- *query, mid chan interface{}, in <-chan *body, qt QueryTranslator) {
	go qt.startParser(mid, in)
	go qt.startBuilder(out, mid)
}

func startExecutor(in <-chan *query) <-chan time.Duration {
	elapsed := make(chan time.Duration)
	go func() {
		db, err := sql.Open(conf.DB.Driver, conf.DB.DNS)
		if err != nil {
			log.Print(err)
		}
		defer func() {
			db.Close()
			close(elapsed)
		}()
		for q := range in {
			start := time.Now()
			_, err := db.Exec(q.sql, q.params...)
			elapsed <- time.Since(start)
			if err != nil {
				log.Print(err)
			}
		}
	}()
	return elapsed
}

func startMonitor(ep1, ep2 <-chan time.Duration) {
	start := time.Now()
	defer func() {
		log.Printf("total elapsed time: %s", time.Since(start))
	}()
	total := time.After(timeoutSec * time.Second)
loop:
	for {
		select {
		case t1, ok := <-ep1:
			if !ok {
				ep1 = nil
			} else {
				log.Printf("fetching elapsed time %s", t1)
			}
		case t2, ok := <-ep2:
			if !ok {
				break loop
			}
			log.Printf("inserting elapsed time %s", t2)
		case <-total:
			log.Printf("timeout after %d sec", timeoutSec)
			break loop
		}
	}
}

// Fetcher is a basic struct for fetching
type Fetcher struct {
	tpl string
	qt  QueryTranslator
}

// StartFetcher will crawl and store tiantian realtime fund data
func (f Fetcher) StartFetcher(codes []string) {
	sent := make(chan string)
	crawled := make(chan *body)
	parsed := make(chan interface{})
	made := make(chan *query)

	n := len(codes)
	wg1.Add(n)
	wg2.Add(n)

	ep1 := startCrawler(crawled, sent)
	go func() {
		defer close(sent)
		for _, code := range codes {
			sent <- fmt.Sprintf(f.tpl, code)
		}
	}()
	startTranslator(made, parsed, crawled, f.qt)
	ep2 := startExecutor(made)

	startMonitor(ep1, ep2)
}

// NewFetcher creates a new Fetcher
func NewFetcher(tpl string, qt QueryTranslator) Fetcher {
	return Fetcher{tpl, qt}
}
