package gameentitytypes

type CoinCollector interface {
	AddCoinCount(amount int)
	CoinCount() int
}
