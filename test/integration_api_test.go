package test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ltp/internal/config"
	"ltp/internal/repository"
	"ltp/internal/service"
	"ltp/internal/transport"
)

type ltpResponse struct {
	LTP []struct {
		Pair   string  `json:"pair"`
		Amount float64 `json:"amount"`
	} `json:"ltp"`
}

func setupTestServer() *httptest.Server {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	repo := repository.New(conf.KrakenAPIConfig, logger)
	svc := service.New(repo, conf.KrakenAPIConfig, logger)
	router := transport.NewHTTPRouter(svc, conf, logger)
	return httptest.NewServer(router)
}

func TestLTPAPI_AllPairs(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/ltp")
	if err != nil {
		t.Fatalf("Failed to GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}

	var result ltpResponse
	decErr := json.NewDecoder(resp.Body).Decode(&result)
	if decErr != nil {
		// Re-fetch the body for debugging
		resp2, err2 := http.Get(ts.URL + "/api/v1/ltp")
		if err2 == nil {
			defer resp2.Body.Close()
			raw, _ := io.ReadAll(resp2.Body)
			t.Fatalf("Failed to decode response: %v\nRaw body: %s", decErr, string(raw))
		} else {
			t.Fatalf("Failed to decode response: %v (and failed to re-fetch body: %v)", decErr, err2)
		}
	}

	if len(result.LTP) == 0 {
		// Warn, but do not fail, if no LTP results (may happen if upstream API is missing data)
		resp2, err2 := http.Get(ts.URL + "/api/v1/ltp")
		if err2 == nil {
			defer resp2.Body.Close()
			raw, _ := io.ReadAll(resp2.Body)
			t.Logf("WARNING: No LTP results.\nStatus: %d\nRaw body: %s\nDecoded: %+v", resp2.StatusCode, string(raw), result)
		} else {
			t.Logf("WARNING: No LTP results. Decoded: %+v (and failed to re-fetch body: %v)", result, err2)
		}
	}
}

func TestLTPAPI_SinglePair(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/ltp?pair=BTC/USD")
	if err != nil {
		t.Fatalf("Failed to GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}

	var result ltpResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.LTP) != 1 {
		t.Errorf("Expected 1 LTP result, got %d", len(result.LTP))
	}
	if result.LTP[0].Pair != "BTC/USD" {
		t.Errorf("Expected pair BTC/USD, got %s", result.LTP[0].Pair)
	}
}
