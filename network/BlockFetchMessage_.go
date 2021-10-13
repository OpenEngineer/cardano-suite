package network

import (
	"github.com/christianschmitz/cardano-suite/codec"
	"github.com/christianschmitz/cardano-suite/ledger"
)

type BlockFetchMessage_ interface {
}

type BlockFetchRequestRange_ struct {
	Lower *ledger.Point
	Upper *ledger.Point
}

type BlockFetchClientDone_ struct {
}

type BlockFetchStartBatch_ struct {
}

type BlockFetchNoBlocks_ struct {
}

type BlockFetchBlock_ struct {
	Block *ledger.Block `wrapped`
}

type BlockFetchBatchDone_ struct {
}

var _blockMessageOk = codec.RegisterInterface(
	(*BlockFetchMessage_)(nil),
	&BlockFetchRequestRange_{},
	&BlockFetchClientDone_{},
	&BlockFetchStartBatch_{},
	&BlockFetchNoBlocks_{},
	&BlockFetchBlock_{},
	&BlockFetchBatchDone_{},
)
