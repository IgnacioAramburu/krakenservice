package repository

type KrakenTickerResponse struct {
	Error  []string                  `json:"error"`
	Result map[string]KrakenPairData `json:"result"`
}

type KrakenPairData struct {
	C []string `json:"c"`
}
