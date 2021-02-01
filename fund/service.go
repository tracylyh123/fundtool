package fund

func GetFundTrend(code, from, to string) (Trend, error) {
	trend, err := FindTrend(code, from, to)
	if err != nil {
		return nil, err
	}
	return trend, nil
}
