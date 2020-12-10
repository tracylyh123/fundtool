package chart

import "github.com/tracylyh123/fundtool/fund"

// List is a set of metric values of trend
type List []float64

// Max returns max item in List
func (l List) Max() float64 {
	if len(l) == 0 {
		panic("the list is empty")
	}
	max := l[0]
	for _, v := range l {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns min item in List
func (l List) Min() float64 {
	if len(l) == 0 {
		panic("the list is empty")
	}
	min := l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return min
}

// CreatePoints creates points which will be plotted
func (l List) CreatePoints(width, height float64) Points {
	var points Points

	if len(l) == 0 {
		return points
	}

	var (
		step = width / float64(len(l)-1)
		max  = l.Max()
		min  = l.Min()
		mid  = (abs(max) + abs(min)) / 2
		low  = min - mid*0.5
		high = max + mid*0.5
		add  = abs(low)
		top  = high + add
	)

	for i := 0; i < len(l); i++ {
		points = append(points, Point{
			x: float64(i) * step,
			y: height * (1 - 1/top*(l[i]+add)),
		})
	}
	return points
}

// NewUnitPriceList creates a new list of unit price
func NewUnitPriceList(t fund.Trend) List {
	var r List
	for _, v := range t {
		r = append(r, float64(v.Price)/10000)
	}
	return r
}

// NewGrowthRateList creates a new list of growth rate
func NewGrowthRateList(t fund.Trend) List {
	r := List{0}
	for i := 1; i < len(t); i++ {
		r = append(r, float64(t[i].Price-t[0].Price)/float64(t[0].Price))
	}
	return r
}
