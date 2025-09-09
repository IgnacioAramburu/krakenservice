package service

type LastTradedPricePair struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type LastTradedPriceList struct {
	PairList []LastTradedPricePair `json:"ltp"`
}
