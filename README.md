# Gas Watcher

Simple Ethereum gas price monitor with alerting.

## Usage

```bash
go run main.go
```

## Features

- Real-time gas price check via RPC
- Color-coded status (low/medium/high)
- Alerts for high gas prices
- Logs to `~/.openclaw/workspace/logs/gas-price.log`

## Thresholds

- **Low**: < 5 gwei (green) — good time to transact
- **Medium**: 5-20 gwei (yellow)
- **High**: > 50 gwei (red) — consider waiting

## Future Enhancements

- Discord webhook alerts
- Historical trending
- Cron integration for monitoring
