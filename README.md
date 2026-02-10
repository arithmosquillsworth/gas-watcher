# Gas Watcher

Simple Ethereum gas price monitor with alerting.

## Usage

```bash
# Basic usage
go run main.go

# With Discord alerts
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..."
go run main.go
```

## Features

- Real-time gas price check via RPC
- Color-coded status (low/medium/high)
- Alerts for high gas prices
- Discord webhook integration
- Logs to `~/.openclaw/workspace/logs/gas-price.log`

## Thresholds

- **Low**: < 5 gwei (green) — good time to transact
- **Medium**: 5-20 gwei (yellow)
- **High**: > 50 gwei (red) — consider waiting

## Environment Variables

- `DISCORD_WEBHOOK_URL` — Optional webhook for alerts

## Future Enhancements

- Historical trending
- Cron integration for monitoring
