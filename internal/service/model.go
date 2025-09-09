package service

import "encoding/json"

type LastTradedPricePair struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type LastTradedPriceList struct {
	PairList []LastTradedPricePair `json:"ltp"`
}

// Ensure non-nil slice when marshaling
func (l *LastTradedPriceList) MarshalJSON() ([]byte, error) {
	type Alias LastTradedPriceList
	if l.PairList == nil {
		l.PairList = []LastTradedPricePair{}
	}
	return json.Marshal((*Alias)(l))
}
