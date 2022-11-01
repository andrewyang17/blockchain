// Package selector provides different transaction selecting algorithms.
package selector

import (
	"fmt"
	"strings"

	"github.com/andrewyang17/blockchain/foundation/blockchain/database"
)

const (
	StrategyTip         = "tip"
	StrategyTipAdvanced = "tip_advanced"
)

var strategies = map[string]Func{
	StrategyTip: tipSelect,
}

// Func defines a function that takes a mempool of transactions grouped by
// account and selects howMany of them in an order based on the function strategy.
// All selector function MUST respect nonce ordering. Receiving for howMany
// must return all the transactions in the strategies ordering.
type Func func(transaction map[database.AccountID][]database.BlockTx, howMany int) []database.BlockTx

// Retrieve returns the specified select strategy function.
func Retrieve(strategy string) (Func, error) {
	fn, exists := strategies[strings.ToLower(strategy)]
	if !exists {
		return nil, fmt.Errorf("strategy %q does not exist", strategy)
	}
	return fn, nil
}

// =============================================================================

// byNonce provides sorting support by the transaction id value.
type byNonce []database.BlockTx

func (b byNonce) Len() int {
	return len(b)
}

// Less helps to sort the list by nonce in ascending order to keep the
// transactions in the right order of processing.
func (b byNonce) Less(i, j int) bool {
	return b[i].Nonce < b[j].Nonce
}

func (b byNonce) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// =============================================================================

// byTip provides sorting support by the transaction tip value.
type byTip []database.BlockTx

func (b byTip) Len() int {
	return len(b)
}

// Less helps to sort the list by tip in decending order to pick the
// transactions that provide the best reward.
func (b byTip) Less(i, j int) bool {
	return b[i].Tip > b[j].Tip
}

func (b byTip) Swap(i, j int) {
	//TODO implement me
	panic("implement me")
}
