package domain

//Money define o tipo do valor monetário
type Money int64

//Float64 converte o tipo Money para float64
func (m Money) Float64() float64 {
	return float64(m) / 100
}
