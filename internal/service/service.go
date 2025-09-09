package service

import (
	"context"
	"fmt"
	"log"
	"ltp/internal/_errors"
	"ltp/internal/config"
	"ltp/internal/repository"
	"strconv"
	"strings"
)

type Service interface {
	GetLastTradedPrice(ctx context.Context, pairInfo string) (ltp_list *LastTradedPriceList, errInfo *_errors.ErrorInformation)
}

type service struct {
	repository   repository.Repository
	krakenConfig config.KrakenAPIConfigInfo
	log          *log.Logger
}

func (s *service) GetLastTradedPrice(ctx context.Context, pairInfo string) (ltp_list *LastTradedPriceList, errInfo *_errors.ErrorInformation) {
	tradeCodes, err := s.getTradeCodes(pairInfo)

	if err != nil {
		s.log.Println(err)
		return nil, &_errors.ErrorInformation{
			Type:        _errors.InvalidInputErrType,
			Description: err.Error(),
		}
	}

	krakenTickerInfo, errInfo := s.repository.GetLastTradePriceInquiry(ctx, tradeCodes)

	if errInfo != nil {
		s.log.Println(errInfo.Description)
		return nil, errInfo
	}

	//Not all pairs are supported by Kraken. So derivative pairs are calculated from the ones supported according
	//The tradeCode expressed as chain <step1>-><step2>->...-><stepn> in s.krakenConfig.PairMap

	s.adaptDerivativePairs(krakenTickerInfo)

	var ltpList []LastTradedPricePair

	for tradeCode, data := range krakenTickerInfo.Result {
		if len(data.C) == 0 {
			continue // skip if no trade data
		}

		priceStr := data.C[0]

		price, err := strconv.ParseFloat(priceStr, 64)

		if err != nil {
			s.log.Printf("Invalid price for %s: %v", tradeCode, err)
			continue
		}

		pair, ok := s.krakenConfig.TradeCodesMap[tradeCode]

		if !ok {
			continue
		}

		if strings.Contains(pairInfo, pair) { //Only add requested pairs
			ltpList = append(ltpList, LastTradedPricePair{
				Pair:   pair,
				Amount: price,
			})
		}
	}

	return &LastTradedPriceList{PairList: ltpList}, nil

}

func (s *service) getTradeCodes(pairInfo string) (string, error) {
	if pairInfo == "" {
		return s.krakenConfig.AvailableTradeCodes, nil
	}
	return s.getSpecificTradeCodes(pairInfo)
}

func (s *service) getSpecificTradeCodes(pairInfo string) (string, error) {
	pairs := strings.Split(pairInfo, ",")
	tradeCodes := ""
	for _, pair := range pairs {
		val, ok := s.krakenConfig.PairMap[pair]
		if !ok {
			err := fmt.Errorf("invalid Trade Pair: %s", pair)
			return "", err
		}
		tradeCodes += val + ","
	}
	return tradeCodes[:len(tradeCodes)-1], nil
}

func (s *service) adaptDerivativePairs(krakenTickerInfo *repository.KrakenTickerResponse) {
	for tradeCode := range s.krakenConfig.TradeCodesMap {
		if strings.Contains(tradeCode, "->") {
			tradeChain := strings.Split(tradeCode, "->")
			finalPrice := 1.0
			finalVolume := 1.0
			valid := true

			for _, step := range tradeChain {
				data, ok := krakenTickerInfo.Result[step]
				if !ok || len(data.C) < 2 {
					valid = false
					break
				}
				price, err1 := strconv.ParseFloat(data.C[0], 64)
				volume, err2 := strconv.ParseFloat(data.C[1], 64)
				if err1 != nil || err2 != nil || price == 0 {
					valid = false
					break
				}
				finalPrice *= price
				finalVolume *= volume
			}

			if valid {
				krakenTickerInfo.Result[tradeCode] = repository.KrakenPairData{
					C: []string{
						fmt.Sprintf("%.8f", finalPrice),
						fmt.Sprintf("%.8f", finalVolume),
					},
				}
				s.log.Printf("Derived synthetic pair %s = [%.8f, %.8f]", tradeCode, finalPrice, finalVolume)
			}
		}
	}
}

func New(repo repository.Repository, krakenConf config.KrakenAPIConfigInfo, logger *log.Logger) Service {
	return &service{
		repository:   repo,
		krakenConfig: krakenConf,
		log:          logger,
	}
}
