# Golang Developer Assigment

Develop in Go language a service that will provide an API for retrieval of the Last Traded Price of Bitcoin for the following currency pairs:

1. BTC/USD
2. BTC/CHF
3. BTC/EUR


The request path is:
/api/v1/ltp

The response shall constitute JSON of the following structure:
```json
{
  "ltp": [
    {
      "pair": "BTC/CHF",
      "amount": 49000.12
    },
    {
      "pair": "BTC/EUR",
      "amount": 50000.12
    },
    {
      "pair": "BTC/USD",
      "amount": 52000.12
    }
  ]
}

```

# Requirements:

1. The incoming request can done for as for a single pair as well for a list of them
2. You shall provide time accuracy of the data up to the last minute.
3. Code shall be hosted in a remote public repository
4. readme.md includes clear steps to build and run the app
5. Integration tests
6. Dockerized application

# Docs
The public Kraken API might be used to retrieve the above LTP information
[API Documentation](https://docs.kraken.com/rest/#tag/Spot-Market-Data/operation/getTickerInformation)
(The values of the last traded price is called “last trade closed”)

# Solution

## Build and Run with Docker

You can build and run the containerized service using Docker and with the provided commands in the Makefile:

(Service is called krakenservice as container name)

### Build the Docker image
```sh
make docker-build
```

### Run the service in Docker
```sh
make docker-run
```
This will start the service and expose it on port 8080.

### Stop the running Docker container
```sh
make docker-stop
```

### Remove the Docker image and local binary
```sh
make docker-clean
```

### Build the Go binary locally (optional)
```sh
make build
```

## How to make API request

Call the API this way:

<domain>/api/v1/ltp?pair=<pairs>

* domain: domain where is hosted. If run locally domain is http://127.0.0.1 (default port is 8080)
* pairs: pair="<pair1>,...,<pairn>",being pair any of these supported pairs: "BTC/USD", "BTC/EUR", "BTC/CHF". You can request one or many pairs at the same time as you wish.

if <pairs> is empty string or no pairs in query params, it will return the last trade price or all supported pairs

## Summary of Makefile Commands

- `make build` – Build the Go binary locally
- `make docker-build` – Build the Docker image
- `make docker-run` – Run the Docker container (exposes port 8080)
- `make docker-stop` – Stop the running container
- `make docker-clean` – Remove the Docker image and local binary
