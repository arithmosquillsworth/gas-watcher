# Gas Watcher

Simple Ethereum gas price monitor with alerting.

## Usage

```bash
# Interactive mode (with UI)
go run main.go

# Quiet/cron mode (logs only, alerts when threshold hit)
go run main.go --quiet

# With Discord alerts
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..."
go run main.go
```

## Cron Setup

Check gas every 15 minutes and alert on Discord:

```bash
*/15 * * * * cd /home/arithmos/projects/arithmosquillsworth/gas-watcher && DISCORD_WEBHOOK_URL="your-webhook-url" ./gas-watcher --quiet
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
