package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tracylyh123/fundtool/fetcher"
)

// EST = fund with estimate value, NET = fund with net value
const (
	EST = iota
	NET
)

// Kind will contain different template and translator according to kind flag
type Kind struct {
	tpl  string
	kind int
	qt   fetcher.QueryTranslator
}

// Set initializes Kind by flag
func (f *Kind) Set(s string) error {
	switch s {
	case "est":
		f.tpl = "http://fundgz.1234567.com.cn/js/%s.js"
		f.qt = &fetcher.FundEstQueryTranslator{}
		f.kind = EST
		return nil
	case "net":
		f.tpl = "http://fund.eastmoney.com/pingzhongdata/%s.js"
		f.qt = &fetcher.FundNetHistoryQueryTranslator{}
		f.kind = NET
		return nil
	}
	return fmt.Errorf("invalid kind %q", s)
}

func (f *Kind) String() string {
	return fmt.Sprintf("tpl: %s, kind: %d", f.tpl, f.kind)
}

// KindFlag defines a Kind flag
func KindFlag(name string, value string, usage string) *Kind {
	f := &Kind{}
	flag.CommandLine.Var(f, name, usage)
	return f
}

// Codes is a set of fund code
type Codes []string

// Set splits string to slice by comma
func (f *Codes) Set(s string) error {
	for _, i := range strings.Split(s, ",") {
		*f = append(*f, i)
	}
	return nil
}

func (f *Codes) String() string {
	return fmt.Sprintf("%v", *f)
}

// CodeFlag defines a Code flag
func CodeFlag(name string, value string, usage string) *Codes {
	f := &Codes{}
	flag.CommandLine.Var(f, name, usage)
	return f
}
