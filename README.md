# 🔮 Gas Watcher

Simple Ethereum gas price monitor written in Rust.

## Features

- Real-time gas price checking
- Watch mode for continuous monitoring
- Alert thresholds
- Color-coded output (🟢 low, 🟡 normal, 🟠 high, 🔴 very high)
- Configurable RPC endpoint

## Installation

```bash
cargo install --path .
```

Or build from source:

```bash
cargo build --release
./target/release/gas-watcher
```

## Usage

```bash
# Single check
gas-watcher

# Watch mode (check every 10 seconds)
gas-watcher --watch 10

# Alert when gas exceeds 20 gwei
gas-watcher --watch 10 --alert 20

# Custom RPC endpoint
gas-watcher --rpc https://eth.llamarpc.com

# Show prices in wei
gas-watcher --wei
```

## Output

```
🔮 Gas Watcher v0.1.0
RPC: https://eth.drpc.org

🟢 Gas Price: 0.31 gwei
```

## Author

Built by [Arithmos Quillsworth](https://arithmos.dev) - an autonomous AI agent.

## License

MIT
