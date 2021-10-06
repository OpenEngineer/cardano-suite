package cardano

type BlockFetchMessage interface {
}

type BlockFetchRequestRange struct {
  Lower *Point
  Upper *Point // inclusive, so for fetching a single block use "BlockFetchRequestRange{PointA, PointA}"
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

//go:generate ./cbor-type BlockFetchMessage *BlockFetchRequestRange *BlockFetchClientDone *BlockFetchStartBatch *BlockFetchNoBlocks *BlockFetchBlock *BlockFetchBatchDone
//go:generate ./cbor-type *BlockFetchRequestRange "Lower *Point, Upper *Point"
//go:generate ./cbor-type *BlockFetchClientDone
//go:generate ./cbor-type *BlockFetchStartBatch
//go:generate ./cbor-type *BlockFetchNoBlocks
//go:generate ./cbor-type *BlockFetchBlock "Block *WrappedBlock"
//go:generate ./cbor-type *BlockFetchBatchDone
