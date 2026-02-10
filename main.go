package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	rpcEndpoint = "https://eth.drpc.org"
	lowThreshold    = 5.0   // gwei
	mediumThreshold = 20.0  // gwei
	highThreshold   = 50.0  // gwei
)

type RPCRequest struct {
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

type RPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

func getGasPrice() (float64, error) {
	reqBody := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_gasPrice",
		Params:  []string{},
		ID:      1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(rpcEndpoint, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return 0, err
	}

	// Convert hex to decimal
	hexValue := strings.TrimPrefix(rpcResp.Result, "0x")
	wei, err := strconv.ParseInt(hexValue, 16, 64)
	if err != nil {
		return 0, err
	}

	// Convert wei to gwei
	gwei := float64(wei) / 1e9
	return gwei, nil
}

func getStatusColor(gwei float64) string {
	if gwei < lowThreshold {
		return "\033[32m" // Green
	} else if gwei < mediumThreshold {
		return "\033[33m" // Yellow
	}
	return "\033[31m" // Red
}

func getStatusText(gwei float64) string {
	if gwei < lowThreshold {
		return "LOW"
	} else if gwei < mediumThreshold {
		return "MEDIUM"
	}
	return "HIGH"
}

func logGasPrice(gwei float64, status string) error {
	logDir := os.Getenv("HOME") + "/.openclaw/workspace/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := logDir + "/gas-price.log"
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s | %.2f gwei | %s\n", timestamp, gwei, status)
	_, err = f.WriteString(line)
	return err
}

func main() {
	reset := "\033[0m"
	
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║              ⛽ GAS PRICE ALERT —", time.Now().Format("15:04"), "              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	gwei, err := getGasPrice()
	if err != nil {
		fmt.Printf("Error fetching gas price: %v\n", err)
		os.Exit(1)
	}

	status := getStatusText(gwei)
	color := getStatusColor(gwei)

	fmt.Printf("Current Gas Price: %s%.2f gwei%s\n", color, gwei, reset)
	fmt.Printf("Status: %s%s%s\n", color, status, reset)
	fmt.Println()
	fmt.Println("Thresholds:")
	fmt.Printf("  Low:    < %.0f gwei\n", lowThreshold)
	fmt.Printf("  Medium: %.0f-%.0f gwei\n", lowThreshold, mediumThreshold)
	fmt.Printf("  High:   > %.0f gwei\n", highThreshold)
	fmt.Println()

	// Alert if high
	if gwei > highThreshold {
		fmt.Printf("%s⚠️  ALERT: High gas prices detected!%s\n", color, reset)
		fmt.Println("   Consider waiting for lower gas before transacting.")
	}

	// Alert if low (good time to transact)
	if gwei < lowThreshold {
		fmt.Printf("%s✓ Good time to transact — gas is low%s\n", color, reset)
	}

	// Log to file
	if err := logGasPrice(gwei, status); err != nil {
		fmt.Printf("Warning: Failed to log: %v\n", err)
	}

	fmt.Println("\nLog saved to: ~/.openclaw/workspace/logs/gas-price.log")
}
