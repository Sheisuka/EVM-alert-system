package evm

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type erc20worker struct {
	logs     chan types.Log
	query    ethereum.FilterQuery
	provider eventsProvider
	consumer eventConsumer
}

func (s *erc20worker) Run(ctx context.Context) {
	s.logs = make(chan types.Log)
	sub, err := s.provider.SubscribeFilterLogs(ctx, s.query, s.logs)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Println("subscription error, reconnecting:", err)
			time.Sleep(5 * time.Second)
		case <-ctx.Done():
			log.Println("worker finished the job")
			return
		case vLog := <-s.logs:
			log.Println(vLog)
		}
	}
}

func (s *erc20worker) Stop() {
	close(s.logs)
}
