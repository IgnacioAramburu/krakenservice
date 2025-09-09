# ---- Build Stage ----
FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o krakenservice ./cmd/main.go

# ---- Run Stage ----
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/krakenservice ./krakenservice
EXPOSE 8080
ENTRYPOINT ["/app/krakenservice"]
