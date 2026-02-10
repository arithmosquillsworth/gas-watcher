package main

import (
	"bytes"
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

type DiscordWebhook struct {
	Content string `json:"content"`
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

func sendDiscordAlert(gwei float64, status string) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return nil // Skip if no webhook configured
	}

	var message string
	if gwei > highThreshold {
		message = fmt.Sprintf("⛽ **High Gas Alert**\nCurrent: %.2f gwei\nStatus: %s\nConsider waiting before transacting.", gwei, status)
	} else if gwei < lowThreshold {
		message = fmt.Sprintf("✅ **Low Gas Opportunity**\nCurrent: %.2f gwei\nGood time to transact!", gwei)
	} else {
		return nil // No alert needed for medium
	}

	payload := DiscordWebhook{Content: message}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("discord webhook returned %d", resp.StatusCode)
	}
	return nil
}

func main() {
	// Check for quiet/cron mode
	quietMode := len(os.Args) > 1 && os.Args[1] == "--quiet"
	
	if !quietMode {
		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║              ⛽ GAS PRICE ALERT —", time.Now().Format("15:04"), "              ║")
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
		fmt.Println()
	}

	gwei, err := getGasPrice()
	if err != nil {
		fmt.Printf("Error fetching gas price: %v\n", err)
		os.Exit(1)
	}

	status := getStatusText(gwei)

	if quietMode {
		// Quiet mode: just log and send Discord alerts
		if err := logGasPrice(gwei, status); err != nil {
			fmt.Printf("Warning: Failed to log: %v\n", err)
		}
		if err := sendDiscordAlert(gwei, status); err != nil {
			fmt.Printf("Warning: Discord alert failed: %v\n", err)
		}
		// Only output if alert-worthy
		if gwei > highThreshold || gwei < lowThreshold {
			fmt.Printf("%.2f gwei (%s)\n", gwei, status)
		}
		return
	}

	// Interactive mode
	color := getStatusColor(gwei)
	reset := "\033[0m"

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

	// Send Discord alert if configured
	if err := sendDiscordAlert(gwei, status); err != nil {
		fmt.Printf("Warning: Discord alert failed: %v\n", err)
	}

	fmt.Println("\nLog saved to: ~/.openclaw/workspace/logs/gas-price.log")
}
