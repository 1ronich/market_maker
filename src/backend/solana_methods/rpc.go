package solana_methods

import (
	"amm/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gagliardetto/solana-go"
)

func GetLatestBlockhash(cluster string) types.BlockhashResponse {
	var rpc_url string = ""
	switch cluster {
	case "devnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_DEV")
	case "mainnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	}
	//rpc_url = "https://api.devnet.solana.com"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getLatestBlockhash",
		"params": map[string]interface{}{
			"commitment":     "processed",
			"minContextSlot": 1000,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending http req: %v\n", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var Response types.BlockhashResponse

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		log.Fatalf("Error umarshalling http req: %v\n", erra)
	}

	return Response
}

func GetTokenAccountBalance(Vault, cluster string) types.GetTokenBalanceResponse {
	var rpc_url string = ""
	switch cluster {
	case "devnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_DEV")
	case "mainnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	}
	//rpc_url = "https://api.devnet.solana.com"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenAccountBalance",
		"params": []interface{}{
			Vault,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending http req: %v\n", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var Response types.GetTokenBalanceResponse

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		log.Fatalf("Error umarshalling http req: %v\n", erra)
	}

	return Response
}

func GetAccountInfo(public_key any, cluster string) types.SolanaAccountResponse {
	var rpc_url string = ""
	switch cluster {
	case "devnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_DEV")
	case "mainnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	}
	//rpc_url = "https://api.devnet.solana.com"

	var pubkey string
	//this if allows both strings and solana.PublicKeys to be passed as parameters
	if reflect.TypeOf(public_key).Kind() == reflect.String {
		pubkey = public_key.(string)
	} else if reflect.TypeOf(public_key).Name() == "PublicKey" {
		pubkey = public_key.(solana.PublicKey).String()
	}

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getAccountInfo",
		"params": []interface{}{
			pubkey,
			map[string]interface{}{
				"encoding":   "base64",
				"commitment": "confirmed",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var Response types.SolanaAccountResponse // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	return Response
}

func GetProgramAccounts(programID string, accountSize uint64) types.GettProgramAccounts { //gets all lps for x program id
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getProgramAccounts",
		"params": []interface{}{
			programID,
			map[string]interface{}{
				"filters": []interface{}{
					map[string]interface{}{
						"dataSize": accountSize,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	rpc_url := os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	//rpc_url = "https://api.devnet.solana.com"

	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var Response types.GettProgramAccounts // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	return Response
}

func GetTokenAccounsByOwner(programID, mint, owner string) types.TokenAccountsResponse { //gets all lps for x program id
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenAccountsByOwner",
		"params": []interface{}{
			owner, // owner
			map[string]interface{}{
				"mint": mint,
			},
			map[string]interface{}{
				"encoding": "jsonParsed",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	rpc_url := os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	//rpc_url = "https://api.devnet.solana.com"

	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
	var Response types.TokenAccountsResponse // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	return Response
}

func GetSignatureStatuses(tx_hash, cluster string) types.SignatureStatusesResponse {
	var rpc_url string = ""
	switch cluster {
	case "devnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_DEV")
	case "mainnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	}
	//rpc_url = "https://api.devnet.solana.com"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getSignatureStatuses",
		"params": []interface{}{
			[]interface{}{
				tx_hash,
			},
			map[string]interface{}{
				"searchTransactionHistory": true,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("raw JSON:", string(body))

	var Response types.SignatureStatusesResponse // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	return Response
}

func GetBalance(account, cluster string) types.GetBalanceResponse {
	var rpc_url string = ""
	switch cluster {
	case "devnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_DEV")
	case "mainnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	}
	//rpc_url = "https://api.devnet.solana.com"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params": []interface{}{
			account,
			map[string]string{
				"commitment": "finalized",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println("raw JSON:", string(body))

	var Response types.GetBalanceResponse // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	return Response
}
