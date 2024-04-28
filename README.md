# prepaidGas Server

The server imcludes `validator` and `executor` scripts

- The `validator` script accepts pending transactions and provides them with confirming signature
- The `executor` script accepts nice orders and plans corresponding transactions for execution

## HTTP

The validator runs http endpoint:

- POST `/validate`
  - requires body `{ "origSign": {hexstr 65 bytes user signature}, "message": { "from": {hexstr 40 bytes address}, "nonce": {hexstr number}, "order": {hexstr number}, "start": {hexstr unix timestamp}, "to": {hexstr 40 bytes address}, "gas": {hexstr number}, "data": {hexstr unlimited size data} } }`
  - returns `{hexstr 65 bytes confirm signature}`
- GET `/load`
  - requires parameters `?offset={number}&reverse={boolean}`
  - returns `{ "id": {uint64 database id}, "validSign": {hexstr 65 bytes confirm signature}, "origSign": {hexstr 65 bytes user signature}, "message": { "from": {hexstr 40 bytes address}, "nonce": {hexstr number}, "order": {hexstr number}, "start": {hexstr unix timestamp}, "to": {hexstr 40 bytes address}, "gas": {hexstr number}, "data": {hexstr unlimited size data} } }`

## DB

To run the scripts running db is required:
`source .env && docker run -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -e POSTGRES_USER=$POSTGRES_USER -p 5432:5432 postgres:15.4`

## Setup

The `.env` file is needed to run the scripts, there is `.env.sample` presented in the repo

Try running with the commands:

```
go run ./cmd/executor/executor.go
go run ./cmd/validator/validator.go
```

## Executor

There is the `IsOrderRisky` function in the `go_modules/utils/order.go` file, modify it to make executor accept more orders
