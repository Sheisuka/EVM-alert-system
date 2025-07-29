package domain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type RuleID uuid.UUID

type RuleKey struct {
	Chain  string
	Type   string // ERC20, eth transfer, etc
	Topics string // serialized [][]common.Hash
}

type Rule struct {
	ID        RuleID
	Key       RuleKey
	Addresses []common.Address
}

type RuleGroup struct {
	Key   RuleKey
	Rules map[RuleID]*Rule
}
