package main

import (
	"context"
	"log"
	"os"

	"github.com/Sheisuka/EVM-alert-system/fetcher/internal/domain"
	"github.com/Sheisuka/EVM-alert-system/fetcher/internal/infra/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file was found")
	}

	endpointURL := os.Getenv("WSS_ADDRESS")
	conn, err := ethclient.Dial(endpointURL)
	if err != nil {
		log.Fatalf("%v", err)
	}

	dummyKey := domain.RuleKey{
		Type: "logs",
	}
	ID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("%v", err)
	}
	dummyRule := &domain.Rule{
		Key:       dummyKey,
		ID:        domain.RuleID(ID),
		Addresses: []common.Address{common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")},
	}
	manager := evm.NewManager(conn, nil)
	manager.Init(context.Background(), []*domain.Rule{dummyRule})
}
