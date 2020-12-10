package fund

const (
	// Spike is a status represents price go up
	Spike = iota
	// Drop is a status represents price go down
	Drop
	// Flat is a status represents price no change
	Flat
)

// Estval is a set of estimate value
type Estval []Price

// Status returns status among Spike, Drop and Flat
func (e Estval) Status() int {
	if len(e) < 1 {
		return Flat
	}
	if e[len(e)-1] > e[0] {
		return Spike
	} else if e[len(e)-1] < e[0] {
		return Drop
	}
	return Flat
}
