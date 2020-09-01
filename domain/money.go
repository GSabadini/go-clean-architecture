package domain

type Money int64

func (m Money) Float64() float64 {
	return float64(m) / 100
}

func (m Money) Int64() int64 {
	return int64(m)
}
