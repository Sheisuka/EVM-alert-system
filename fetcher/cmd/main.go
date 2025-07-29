package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const (
	key = "fetcher/cmd/wallets/UTC--2025-07-23T20-02-09.045524000Z--b200b7ec5146c2b2ace79123cbc80d98669b1e8c"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file was found")
	}

	endpointURL := os.Getenv("WSS_ADDRESS")

	client, err := ethclient.Dial(endpointURL)
	if err != nil {
		log.Fatalf("Failed to acq connection")
	}

	contractAddress := common.HexToAddress("0xCf5540fFFCdC3d510B18bFcA6d2b9987b0772559")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(23022736),
		ToBlock:   big.NewInt(23023136),
		Addresses: []common.Address{contractAddress},
	}

	// logs := make(chan types.Log)
	// sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// for {
	// 	select {
	// 	case err := <-sub.Err():
	// 		log.Fatal(err)
	// 	case vlog := <-logs:
	// 		log.Println(vlog)
	// 	}
	// }

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(logs)
}
