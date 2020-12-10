package chart

// Point is position of coord
type Point struct {
	x, y float64
}

// Points is a set of Point
type Points []Point

func (points Points) connect(plot func(x, y float64)) {
	for i := 1; i < len(points); i++ {
		line(plot, points[i-1].x, points[i-1].y, points[i].x, points[i].y)
	}
}

func line(plot func(x, y float64), x0, y0, x1, y1 float64) {
	steep := abs(y1-y0) > abs(x1-x0)
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := abs(y1 - y0)
	bias := 0.0
	slope := dy / dx
	var ystep float64
	y := y0
	if y0 < y1 {
		ystep = 1
	} else {
		ystep = -1
	}
	for x := x0; x <= x1; x++ {
		if steep {
			plot(y, x)
		} else {
			plot(x, y)
		}
		bias += slope
		if bias >= 0.5 {
			y += ystep
			bias--
		}
	}
}
