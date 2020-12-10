package fund

// Fund is the basic info of a fund
type Fund struct {
	code   string
	netval Netval
	estval Estval
	trend  Trend
}

// Code returns code of fund
func (f Fund) Code() string {
	return f.code
}

// Estval returns estval of fund
func (f Fund) Estval() Estval {
	return f.estval
}

// Trend returns trend of fund
func (f Fund) Trend() Trend {
	return f.trend
}

// IsTradable check if can trade at the moment
func (f Fund) IsTradable(t Moment) bool {
	if f.netval.Price <= 0 {
		return false
	}
	return f.netval.Date.IsSame(t)
}
