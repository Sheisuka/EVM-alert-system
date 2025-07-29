package evm

import (
	"context"

	"github.com/Sheisuka/EVM-alert-system/fetcher/internal/domain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type worker interface {
	Run(context.Context)
	Stop()
}

func NewWorker(provider eventsProvider, consumer eventConsumer, ruleType string, g *domain.RuleGroup) worker {
	if ruleType == "logs" {
		addresses := make(map[common.Address]struct{})
		for _, r := range g.Rules {
			for _, address := range r.Addresses {
				addresses[address] = struct{}{}
			}
		}
		query := ethereum.FilterQuery{
			Addresses: []common.Address{},
		}
		for address := range addresses {
			query.Addresses = append(query.Addresses, address)
		}

		return &erc20worker{
			query:    query,
			provider: provider,
			consumer: consumer,
		}
	}
	return nil
}
