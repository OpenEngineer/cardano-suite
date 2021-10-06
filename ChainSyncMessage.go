package cardano

type ChainSyncMessage interface {
}

type ChainSyncRequestNext struct {
}

type ChainSyncAwaitReply struct {
}

type ChainSyncRollForward struct {
  Header *WrappedBlockHeader // this is serialized twice!
  Tip    *Tip
}

type ChainSyncRollBackward struct {
  Point *Point
  Tip   *Tip
}

type ChainSyncFindIntersect struct {
  Points []*Point
}

type ChainSyncIntersectFound struct {
  Point *Point
  Tip   *Tip
}

type ChainSyncIntersectNotFound struct {
  Tip *Tip
}

type ChainSyncDone struct {
}

//go:generate ./cbor-type ChainSyncMessage *ChainSyncRequestNext *ChainSyncAwaitReply *ChainSyncRollForward *ChainSyncRollBackward *ChainSyncFindIntersect *ChainSyncIntersectFound *ChainSyncIntersectNotFound *ChainSyncDone

//go:generate ./cbor-type *ChainSyncRequestNext
//go:generate ./cbor-type *ChainSyncAwaitReply
//go:generate ./cbor-type *ChainSyncRollForward "Header *WrappedBlockHeader, Tip *Tip"
//go:generate ./cbor-type *ChainSyncRollBackward "Point *Point, Tip *Tip"
//go:generate ./cbor-type *ChainSyncFindIntersect "Points []*Point"
//go:generate ./cbor-type *ChainSyncIntersectFound "Point *Point, Tip *Tip"
//go:generate ./cbor-type *ChainSyncIntersectNotFound "Tip *Tip"
//go:generate ./cbor-type *ChainSyncDone
