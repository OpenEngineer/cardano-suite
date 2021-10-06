package network

import (
  "github.com/christianschmitz/cardano-suite/ledger"
)

type ChainSyncMessage interface {
}

type ChainSyncRequestNext struct {
}

type ChainSyncAwaitReply struct {
}

type ChainSyncRollForward struct {
  Header *WrappedBlockHeader // this is serialized twice!
  Tip    *ledger.Tip
}

type ChainSyncRollBackward struct {
  Point *ledger.Point
  Tip   *ledger.Tip
}

type ChainSyncFindIntersect struct {
  Points []*ledger.Point
}

type ChainSyncIntersectFound struct {
  Point *ledger.Point
  Tip   *ledger.Tip
}

type ChainSyncIntersectNotFound struct {
  Tip *ledger.Tip
}

type ChainSyncDone struct {
}

//go:generate ../cbor-type ChainSyncMessage *ChainSyncRequestNext *ChainSyncAwaitReply *ChainSyncRollForward *ChainSyncRollBackward *ChainSyncFindIntersect *ChainSyncIntersectFound *ChainSyncIntersectNotFound *ChainSyncDone

//go:generate ../cbor-type *ChainSyncRequestNext
//go:generate ../cbor-type *ChainSyncAwaitReply
//go:generate ../cbor-type *ChainSyncRollForward "Header *WrappedBlockHeader, Tip *ledger.Tip"
//go:generate ../cbor-type *ChainSyncRollBackward "Point *ledger.Point, Tip *ledger.Tip"
//go:generate ../cbor-type *ChainSyncFindIntersect "Points []*ledger.Point"
//go:generate ../cbor-type *ChainSyncIntersectFound "Point *ledger.Point, Tip *ledger.Tip"
//go:generate ../cbor-type *ChainSyncIntersectNotFound "Tip *ledger.Tip"
//go:generate ../cbor-type *ChainSyncDone
