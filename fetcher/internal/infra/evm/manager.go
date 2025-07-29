package evm

import (
	"context"
	"errors"
	"sync"

	"github.com/Sheisuka/EVM-alert-system/fetcher/internal/domain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type eventConsumer interface{}

type eventsProvider interface {
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

type manager struct {
	provider eventsProvider
	consumer eventConsumer

	groups  map[domain.RuleKey]*domain.RuleGroup
	workers map[domain.RuleKey]worker

	mu sync.Mutex
}

func NewManager(provider eventsProvider, consumer eventConsumer) *manager {
	return &manager{
		provider: provider,
		consumer: consumer,
		groups:   make(map[domain.RuleKey]*domain.RuleGroup),
		workers:  make(map[domain.RuleKey]worker),
	}
}

func (m *manager) Add(rule *domain.Rule) {
	m.mu.Lock()
	defer m.mu.Unlock()

	g, ok := m.groups[rule.Key]
	if !ok {
		g = &domain.RuleGroup{
			Key:   rule.Key,
			Rules: make(map[domain.RuleID]*domain.Rule),
		}
		m.groups[rule.Key] = g
	}
	g.Rules[rule.ID] = rule
}

func (m *manager) Delete(ruleID domain.RuleID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, g := range m.groups {
		if g.Key != key {
			continue
		}
		if _, exists := g.Rules[ruleID]; !exists {
			return errors.New("rule not found")
		}

		delete(g.Rules, ruleID)
		if len(g.Rules) != 0 {
			// peresborka, return
		}

		m.workers[key].Stop()
		delete(m.workers, key)
		delete(m.groups, key)
	}

	return nil
}

func (m *manager) Init(ctx context.Context, rules []*domain.Rule) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, rule := range rules {
		g, exists := m.groups[rule.Key]
		if !exists {
			g = &domain.RuleGroup{
				Key:   rule.Key,
				Rules: make(map[domain.RuleID]*domain.Rule),
			}
			m.groups[rule.Key] = g
		}
		g.Rules[rule.ID] = rule
	}

	for key, g := range m.groups {
		w := NewWorker(m.provider, m.consumer, "logs", g)
		m.workers[key] = w
		//TODO
		w.Run(ctx)
	}
}
