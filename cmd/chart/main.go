package main

import (
	"log"
	"os"
	"strconv"

	"github.com/tracylyh123/fundtool/chart"
)

const (
	width  = 600
	height = 200
)

func main() {
	var list chart.List
	for _, s := range os.Args[1:] {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Printf("unexpected fund code: %s", s)
			continue
		}
		list = append(list, val)
	}
	m := chart.NewGifMaker(width, height)
	m.Make(list, os.Stdout)
}
