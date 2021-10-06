package ledger

// Point along with blockId
type Tip struct {
  Point   *Point
  BlockId uint64
}

//go:generate ../cbor-type *Tip "Point *Point, BlockId common.Uint64"
