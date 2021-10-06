package ledger

import (
  "errors"

  cbor "github.com/fxamacker/cbor/v2"

  "github.com/christianschmitz/cardano-suite/common"
)
    
type BlockHeader struct {
  HeaderHash  common.Hash // hash([prevBlock, slotNumber, blockNumber, bodyHash])
  PrevBlock   ChainHash
  SlotId      uint64
  BlockId     uint64
  BodyHash    common.Hash
}

//go:generate ../cbor-type *BlockHeader "HeaderHash common.Hash, PrevBlock ChainHash, SlotId common.Uint64, BlockId common.Uint64, BodyHash common.Hash"

func NewBlockHeader(prevBlock ChainHash, slotId uint64, blockId uint64, bodyHash common.Hash) BlockHeader {
  bl := BlockHeader{
    common.NilHash(),
    prevBlock,
    slotId,
    blockId,
    bodyHash,
  }

  hh := bl.ToHash()

  bl.HeaderHash = hh

  return bl
}

func NewVerifiedBlockHeader(headerHash common.Hash, prevBlock ChainHash, slotId uint64, blockId uint64, bodyHash common.Hash) (BlockHeader, error) {

  bl := NewBlockHeader(prevBlock, slotId, blockId, bodyHash)

  if !bl.HeaderHash.Eq(headerHash) {
    return bl, errors.New("wrong headerhash")
  }

  return bl, nil
}

func (h BlockHeader) VerifyHash() error {
  hh := h.ToHash()

  if !h.HeaderHash.Eq(hh) {
    return errors.New("wrong headerhash")
  }

  return nil
}

// recalculates the headerhash
func (bl BlockHeader) ToHash() common.Hash {
  cborFields := make([]interface{}, 4)

  cborFields[0] = bl.PrevBlock.ToUntyped()
  cborFields[1] = bl.SlotId
  cborFields[2] = bl.BlockId
  cborFields[3] = bl.BodyHash.ToUntyped()

  b, err := cbor.Marshal(cborFields)
  if err != nil {
    panic(err)
  }

  return common.HashFromBytes(b)
}
