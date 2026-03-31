# BIKA Game Bot — Go (Mongo-compatible starter)

This project is a Go rewrite starter that keeps using the **same MongoDB database and collection names** as the old Python bot.

## Same Mongo layout
- DB: `bika_slot`
- collections:
  - `users`
  - `groups`
  - `config`
  - `transactions`
  - `orders`

## Included
- `/start`
- `/ping`
- auto group approval request to owner DM
- approve / reject callback buttons
- `/pendinggroups`
- `/groupstatus`
- `.slot 100` with wallet / treasury flow

## Important
This is a structured starter. Depending on exact library versions, a few Telegram field names may need tiny compile-time adjustments.

## Run
```bash
cp .env.example .env
go mod tidy
go run ./cmd/bikagame
```
