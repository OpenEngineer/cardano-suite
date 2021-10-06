package network

import (
  "github.com/christianschmitz/cardano-suite/ledger"
)

type BlockFetchMessage interface {
}

type BlockFetchRequestRange struct {
  Lower *ledger.Point
  Upper *ledger.Point // inclusive, so for fetching a single block use "BlockFetchRequestRange{PointA, PointA}"
}

type BlockFetchClientDone struct {
}

type BlockFetchStartBatch struct {
}

type BlockFetchNoBlocks struct {
}

type BlockFetchBlock struct {
  Block *WrappedBlock
}

type BlockFetchBatchDone struct {
}

//go:generate ../cbor-type BlockFetchMessage *BlockFetchRequestRange *BlockFetchClientDone *BlockFetchStartBatch *BlockFetchNoBlocks *BlockFetchBlock *BlockFetchBatchDone
//go:generate ../cbor-type *BlockFetchRequestRange "Lower *ledger.Point, Upper *ledger.Point"
//go:generate ../cbor-type *BlockFetchClientDone
//go:generate ../cbor-type *BlockFetchStartBatch
//go:generate ../cbor-type *BlockFetchNoBlocks
//go:generate ../cbor-type *BlockFetchBlock "Block *WrappedBlock"
//go:generate ../cbor-type *BlockFetchBatchDone
