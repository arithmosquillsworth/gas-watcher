use clap::Parser;
use serde::{Deserialize, Serialize};
use std::time::Duration;

#[derive(Parser, Debug)]
#[command(name = "gas-watcher")]
#[command(about = "Monitor Ethereum gas prices", long_about = None)]
struct Args {
    /// RPC endpoint URL
    #[arg(short, long, default_value = "https://eth.drpc.org")]
    rpc: String,

    /// Watch mode - poll every N seconds
    #[arg(short, long)]
    watch: Option<u64>,

    /// Alert threshold in gwei
    #[arg(short, long)]
    alert: Option<f64>,

    /// Show prices in wei instead of gwei
    #[arg(long)]
    wei: bool,
}

#[derive(Serialize)]
struct RpcRequest {
    jsonrpc: &'static str,
    method: &'static str,
    params: Vec<String>,
    id: u32,
}

#[derive(Deserialize, Debug)]
struct RpcResponse {
    result: Option<String>,
    error: Option<RpcError>,
}

#[derive(Deserialize, Debug)]
struct RpcError {
    message: String,
}

fn get_gas_price(rpc_url: &str) -> Result<u128, Box<dyn std::error::Error>> {
    let request = RpcRequest {
        jsonrpc: "2.0",
        method: "eth_gasPrice",
        params: vec![],
        id: 1,
    };

    let response: RpcResponse = ureq::post(rpc_url)
        .timeout(Duration::from_secs(10))
        .send_json(&request)?
        .into_json()?;

    if let Some(error) = response.error {
        return Err(format!("RPC error: {}", error.message).into());
    }

    let hex_price = response.result.ok_or("No result in response")?;
    let price = u128::from_str_radix(hex_price.trim_start_matches("0x"), 16)?;
    Ok(price)
}

fn wei_to_gwei(wei: u128) -> f64 {
    wei as f64 / 1_000_000_000.0
}

fn format_price(wei: u128, as_wei: bool) -> String {
    if as_wei {
        format!("{} wei", wei)
    } else {
        format!("{:.2} gwei", wei_to_gwei(wei))
    }
}

fn main() {
    let args = Args::parse();

    println!("🔮 Gas Watcher v0.1.0");
    println!("RPC: {}\n", args.rpc);

    loop {
        match get_gas_price(&args.rpc) {
            Ok(price) => {
                let formatted = format_price(price, args.wei);
                let gwei = wei_to_gwei(price);
                
                // Color coding based on price
                let indicator = if gwei < 10.0 {
                    "🟢" // Very low
                } else if gwei < 30.0 {
                    "🟡" // Normal
                } else if gwei < 100.0 {
                    "🟠" // High
                } else {
                    "🔴" // Very high
                };

                println!("{} Gas Price: {}", indicator, formatted);

                // Alert if threshold exceeded
                if let Some(threshold) = args.alert {
                    if gwei > threshold {
                        println!(
                            "⚠️  ALERT: Gas price ({:.2} gwei) exceeds threshold ({:.2} gwei)!",
                            gwei, threshold
                        );
                    }
                }
            }
            Err(e) => {
                eprintln!("❌ Error fetching gas price: {}", e);
            }
        }

        // If not in watch mode, exit after one check
        match args.watch {
            Some(seconds) => {
                std::thread::sleep(Duration::from_secs(seconds));
            }
            None => break,
        }
    }
}
