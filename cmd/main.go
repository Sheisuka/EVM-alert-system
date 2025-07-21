package main

import (
	"alertsys/weth"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const (
	alchemy_api_http = "https://eth-mainnet.g.alchemy.com/v2"
	alchemy_api_wss  = "wss://eth-mainnet.g.alchemy.com/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file was found")
	}

	endpointURL := fmt.Sprintf("%s/%s", alchemy_api_wss, os.Getenv("ALCHEMY_API_KEY"))
	client, err := ethclient.Dial(endpointURL)
	if err != nil {
		log.Fatalf("Failed to acq connection")
	}
	_ = client

	contractAbi, err := abi.JSON((strings.NewReader(string(weth.WethABI))))
	if err != nil {
		log.Fatalf("%v", err)
	}
	_ = contractAbi

	wethAddress := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(22969611),
		ToBlock:   big.NewInt(22969911),
		Addresses: []common.Address{wethAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, vLog := range logs {
		i, err := contractAbi.Unpack("Withdrawal", vLog.Data)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		if v, ok := i[0].(*big.Int); ok {
			log.Println(*v)
			break
		}
	}
}
