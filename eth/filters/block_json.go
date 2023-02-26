package filters

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/params"
)

// Encodes the block to RPC output including full txs.
func marshalBlock(b *types.Block) (map[string]interface{}, error) {
	cfg := new(params.ChainConfig)

	txs := b.Transactions()
	if len(txs) != 0 {
		cfg.ChainID = txs[0].ChainId()
	}

	return ethapi.RPCMarshalBlock(b, true, true, cfg)
}
