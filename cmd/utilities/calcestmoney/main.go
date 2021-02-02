package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"text/tabwriter"
)

const srctpl = "http://fundgz.1234567.com.cn/js/%s.js"

type srcitem struct {
	Code     string `json:"fundcode"`
	Name     string `json:"name"`
	Price    string `json:"dwjz"`
	Estprice string `json:"gsz"`
	Esttime  string `json:"gztime"`
	Estrate  string `json:"gszzl"`
}

type confitem struct {
	Code  string  `json:"code"`
	Money float64 `json:"money"`
}

type estitem struct {
	code     string
	name     string
	money    float64
	price    float64
	estprice float64
	estrate  string
	esttime  string
}

func (e *estitem) update(wg *sync.WaitGroup) {
	url := fmt.Sprintf(srctpl, e.code)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		resp.Body.Close()
		wg.Done()
	}()
	if resp.StatusCode != 200 {
		log.Fatal(fmt.Errorf("invalid http status code %d", resp.StatusCode))
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	f := &srcitem{}
	err = json.Unmarshal([]byte(b[8:len(b)-2]), f)
	if err != nil {
		log.Fatal(err)
	}
	ep, err := strconv.ParseFloat(f.Estprice, 64)
	if err != nil {
		log.Fatal(err)
	}
	e.estprice = ep
	p, err := strconv.ParseFloat(f.Price, 64)
	if err != nil {
		log.Fatal(err)
	}
	e.price = p
	e.name = f.Name
	e.estrate = f.Estrate
	e.esttime = f.Esttime
}

func (e estitem) estmoney() float64 {
	return e.money * e.estprice / e.price
}

func (e estitem) earning() float64 {
	return e.estmoney() - e.money
}

func (e estitem) String() string {
	return fmt.Sprintf("%s current money: %.2f, estimate money: %.2f, estimate rate: %s, earning: %.2f", e.name, e.money, e.estmoney(), e.estrate, e.earning())
}

type estset []*estitem

func (es estset) Len() int           { return len(es) }
func (es estset) Less(i, j int) bool { return es[i].earning() < es[j].earning() }
func (es estset) Swap(i, j int)      { es[i], es[j] = es[j], es[i] }

func (es estset) total() (money, estmoeny, earning, estrate float64) {
	for _, e := range es {
		money += e.money
		estmoeny += e.estmoney()
	}
	earning = estmoeny - money
	estrate = earning / money * 100
	return
}

func (es estset) String() string {
	money, estmoney, earning, estrate := es.total()
	return fmt.Sprintf("Total\ncurrent money: %.2f, estimate money: %.2f, estimate rate: %.2f, earning: %.2f", money, estmoney, estrate, earning)
}

func (es estset) print() {
	sort.Sort(es)
	for _, est := range es {
		fmt.Println(est)
	}
	fmt.Println(es)
}

func (es estset) printf() {
	sort.Sort(es)
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t\n", "Money", "Estimate Money", "Estimate Rate", "Earning", "Name")
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t\n", "-----", "--------------", "-------------", "-------", "----")
	for _, est := range es {
		fmt.Fprintf(tw, "%.2f\t%.2f\t%s\t%.2f\t%s\t\n", est.money, est.estmoney(), est.estrate, est.earning(), est.name)
	}
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t\n", "-----", "--------------", "-------------", "-------", "----")
	money, estmoney, earning, estrate := es.total()
	fmt.Fprintf(tw, "%.2f\t%.2f\t%.2f\t%.2f\t%s\t\n", money, estmoney, estrate, earning, "Total")
	tw.Flush()
}

func main() {
	buffer, err := ioutil.ReadFile("./funds.json")
	if err != nil {
		log.Fatal(err)
	}
	funds := []confitem{}
	err = json.Unmarshal(buffer, &funds)
	if err != nil {
		log.Fatal(err)
	}
	myset := estset{}
	var wg sync.WaitGroup
	for _, fund := range funds {
		wg.Add(1)
		e := estitem{code: fund.Code, money: fund.Money}
		go e.update(&wg)
		myset = append(myset, &e)
	}
	wg.Wait()
	myset.printf()
}
