package network

import (
	"github.com/christianschmitz/cardano-suite/codec"
	"github.com/christianschmitz/cardano-suite/ledger"
)

type ChainSyncMessage_ interface {
}

type ChainSyncRequestNext_ struct {
}

type ChainSyncAwaitReply_ struct {
}

type ChainSyncRollForward_ struct {
	Header *ledger.BlockHeader `wrapped`
	Tip    *ledger.Tip
}

type ChainSyncRollBackward_ struct {
	Point *ledger.Point
	Tip   *ledger.Tip
}

type ChainSyncFindIntersect_ struct {
	Points []*ledger.Point
}

type ChainSyncIntersectFound_ struct {
	Point *ledger.Point
	Tip   *ledger.Tip
}

type ChainSyncIntersectNotFound_ struct {
	Tip *ledger.Tip
}

type ChainSyncDone_ struct {
}

var _chainSyncMessageOk = codec.RegisterInterface(
	(*ChainSyncMessage_)(nil),
	&ChainSyncRequestNext_{},
	&ChainSyncAwaitReply_{},
	&ChainSyncRollForward_{},
	&ChainSyncRollBackward_{},
	&ChainSyncFindIntersect_{},
	&ChainSyncIntersectFound_{},
	&ChainSyncIntersectNotFound_{},
	&ChainSyncDone_{},
)
