# BIKA Game Bot (Go)

Mongo-compatible Go version for the earlier Python bot.

## Features included
- /start
- /ping
- /status
- /balance /bal .bal
- /dailyclaim
- /top10 .top10
- /gift .gift
- /pendinggroups
- /groupstatus
- /approve /reject
- owner DM callback approve/reject
- .slot

## Run
```bash
cp .env.example .env
# edit .env

go mod tidy
go build -o bikagame ./cmd/bikagame
./bikagame
```

## Mongo compatibility
Uses the same database/collections as the Python bot:
- DB: `bika_slot`
- collections: `users`, `groups`, `config`, `transactions`, `orders`
