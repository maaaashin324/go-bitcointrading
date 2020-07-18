# Gotrading

This repository can be used to get bitcoin trading data from bitflyer on websocket.

## How to run this server

### Clone this repository

```bash
git clone 
```

### Get necessary packages

```bash
go get -u
```

### Run server

```bash
go run main.go
```

## Endpoint

This repository runs server that has two endpoints.

### Show chart of candles

Using templates, you can see candles in `/charts/`.

### API Endpoint

You can fetch data from `/api/candle/`.
