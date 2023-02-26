package types

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Decoding block from json. The code is copied from
// https://github.com/bountyHntr/geth/blob/release/1.11/ethclient/ethclient.go#L110
// and modified - unnecessary code sections are commented out.

type rpcBlock struct {
	Hash         common.Hash      `json:"hash"`
	Transactions []rpcTransaction `json:"transactions"`
	UncleHashes  []common.Hash    `json:"uncles"`
}

type rpcTransaction struct {
	tx *Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

func (b *Block) UnmarshalJSON(raw []byte) error {
	// Decode header and transactions.
	var head *Header
	var body rpcBlock
	if err := json.Unmarshal(raw, &head); err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return err
	}
	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	if head.UncleHash == EmptyUncleHash && len(body.UncleHashes) > 0 {
		return fmt.Errorf("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != EmptyUncleHash && len(body.UncleHashes) == 0 {
		return fmt.Errorf("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == EmptyRootHash && len(body.Transactions) > 0 {
		return fmt.Errorf("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != EmptyRootHash && len(body.Transactions) == 0 {
		return fmt.Errorf("server returned empty transaction list but block header indicates transactions")
	}

	// Skipped
	// Load uncles because they are not included in the block response.
	// var uncles []*types.Header
	// if len(body.UncleHashes) > 0 {
	// 	uncles = make([]*types.Header, len(body.UncleHashes))
	// 	reqs := make([]rpc.BatchElem, len(body.UncleHashes))
	// 	for i := range reqs {
	// 		reqs[i] = rpc.BatchElem{
	// 			Method: "eth_getUncleByBlockHashAndIndex",
	// 			Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
	// 			Result: &uncles[i],
	// 		}
	// 	}
	// 	if err := ec.c.BatchCallContext(ctx, reqs); err != nil {
	// 		return nil, err
	// 	}
	// 	for i := range reqs {
	// 		if reqs[i].Error != nil {
	// 			return nil, reqs[i].Error
	// 		}
	// 		if uncles[i] == nil {
	// 			return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
	// 		}
	// 	}
	// }

	// Fill the sender cache of transactions in the block.
	txs := make([]*Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		// Skipped
		// if tx.From != nil {
		// 	setSenderFromServer(tx.tx, *tx.From, body.Hash)
		// }

		txs[i] = tx.tx
	}

	*b = *NewBlockWithHeader(head).WithBody(txs, nil)
	return nil
}
