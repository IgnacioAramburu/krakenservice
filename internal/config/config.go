package config

import (
	"errors"
	"fmt"
	"log"
	"ltp/utils"
)

//Server Config

const ServerPort = "8080"

// Currencies
const (
	CHF = "CHF"
	USD = "USD"
	BTC = "BTC"
	EUR = "EUR"
)

// Valid Pairs
var (
	BTC_USD = fmt.Sprintf("%s/%s", BTC, USD)
	BTC_EUR = fmt.Sprintf("%s/%s", BTC, EUR)
	BTC_CHF = fmt.Sprintf("%s/%s", BTC, CHF)
)

//Kraken Integration

const KrakenUrl = "https://api.kraken.com/0/public/Ticker"

var AvailableKrakenTradePairMap = map[string]string{
	BTC_USD: "XXBTZUSD",
	BTC_EUR: "XXBTZEUR",
	BTC_CHF: "XXBTZUSD->USDCHF", // pair interpreted as trainChain, were each step is delimited by '->' and represent all steps to get pair trade price
}

type KrakenAPIConfigInfo struct {
	ApiUrl              string
	PairMap             map[string]string
	TradeCodesMap       map[string]string
	AvailableTradeCodes string
}

type Config struct {
	KrakenAPIConfig KrakenAPIConfigInfo
	ServerPort      string
}

func getAllAvailableTradeCodes() (string, error) {
	tradeCodes := ""
	if len(AvailableKrakenTradePairMap) == 0 {
		return "", errors.New("missing Available Trade Pairs")
	}
	for _, code := range AvailableKrakenTradePairMap {
		tradeCodes += fmt.Sprintf("%s,", code)
	}
	tradeCodes = tradeCodes[:len(tradeCodes)-1]
	return tradeCodes, nil
}

func New() (Config, error) {
	tradeCodes, err := getAllAvailableTradeCodes()
	if err != nil {
		log.Fatalf("Error at instantiating config: %s", err.Error())
		return Config{}, err
	}

	krakenAPIConfig := KrakenAPIConfigInfo{
		ApiUrl:              KrakenUrl,
		PairMap:             AvailableKrakenTradePairMap,
		TradeCodesMap:       utils.ReverseMap(AvailableKrakenTradePairMap),
		AvailableTradeCodes: tradeCodes,
		//Other required configurations requested by provider should go here
	}

	return Config{
		KrakenAPIConfig: krakenAPIConfig,
		ServerPort:      ServerPort,
	}, nil
}
