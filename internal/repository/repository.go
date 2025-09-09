package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"ltp/internal/_errors"
	"ltp/internal/config"
	"net/http"
	"net/url"
	"strings"
)

type Repository interface {
	GetLastTradePriceInquiry(ctx context.Context, pairs string) (*KrakenTickerResponse, *_errors.ErrorInformation)
}

type repository struct {
	krakenAPIConfig config.KrakenAPIConfigInfo
	log             *log.Logger
}

func (r *repository) GetLastTradePriceInquiry(ctx context.Context, pairs string) (*KrakenTickerResponse, *_errors.ErrorInformation) {
	pairs = normalizePairs(pairs)

	u, _ := url.Parse(r.krakenAPIConfig.ApiUrl)

	q := u.Query()
	q.Set("pair", pairs)
	u.RawQuery = q.Encode()

	log.Println(u.String())

	resp, err := http.Get(u.String())

	if err != nil {
		log.Println(err)
		return nil, &_errors.ErrorInformation{
			Type:        _errors.GatewayErrType,
			Description: fmt.Sprintf("kraken request failed: %v", err),
		}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return nil, &_errors.ErrorInformation{
			Type:        _errors.ResponseReadingErrType,
			Description: fmt.Sprintf("failed to read response: %v", err),
		}
	}

	var krakenResp *KrakenTickerResponse
	if err := json.Unmarshal(body, &krakenResp); err != nil {
		log.Println(err)
		return nil, &_errors.ErrorInformation{
			Type:        _errors.ResponseReadingErrType,
			Description: fmt.Sprintf("failed to parse JSON: %v", err),
		}
	}

	return krakenResp, nil
}

func normalizePairs(pairs string) string {
	pairs = strings.ReplaceAll(pairs, "->", ",")
	pairList := strings.Split(pairs, ",")

	seen := make(map[string]bool)

	var uniquePairList []string
	for _, str := range pairList {
		if !seen[str] {
			seen[str] = true
			uniquePairList = append(uniquePairList, str)
		}
	}

	return strings.Join(uniquePairList, ",")
}

func New(krakenAPIConfig config.KrakenAPIConfigInfo, logger *log.Logger) Repository {
	return &repository{
		krakenAPIConfig: krakenAPIConfig,
		log:             logger,
	}
}
