package main

import (
	"flag"

	"github.com/tracylyh123/fundtool/fetcher"
)

var kind = KindFlag("kind", "est", "kind of fund data, only supports net and est")
var codes = CodeFlag("codes", "", "fund code, if multiple splited by comma")

func main() {
	flag.Parse()
	fetcher.NewFetcher(kind.tpl, kind.qt).StartFetcher(*codes)
}
