package cardano

import (
  "errors"

  cbor "github.com/fxamacker/cbor/v2"
)
    
type BlockHeader struct {
  HeaderHash  Hash // hash([prevBlock, slotNumber, blockNumber, bodyHash])
  PrevBlock   ChainHash
  SlotId      uint64
  BlockId     uint64
  BodyHash    Hash
}

//go:generate ./cbor-type *BlockHeader "HeaderHash Hash, PrevBlock ChainHash, SlotId Uint64, BlockId Uint64, BodyHash Hash"

type WrappedBlockHeader struct {
  Header *BlockHeader
}

func NewBlockHeader(prevBlock ChainHash, slotId uint64, blockId uint64, bodyHash Hash) BlockHeader {
  bl := BlockHeader{
    NilHash(),
    prevBlock,
    slotId,
    blockId,
    bodyHash,
  }

  hh := bl.ToHash()

  bl.HeaderHash = hh

  return bl
}

func NewVerifiedBlockHeader(headerHash Hash, prevBlock ChainHash, slotId uint64, blockId uint64, bodyHash Hash) (BlockHeader, error) {

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
func (bl BlockHeader) ToHash() Hash {
  cborFields := make([]interface{}, 4)

  cborFields[0] = bl.PrevBlock.ToUntyped()
  cborFields[1] = bl.SlotId
  cborFields[2] = bl.BlockId
  cborFields[3] = bl.BodyHash.ToUntyped()

  b, err := cbor.Marshal(cborFields)
  if err != nil {
    panic(err)
  }

  return HashFromBytes(b)
}

func WrappedBlockHeaderFromUntyped(fields interface{}) (*WrappedBlockHeader, error) {
  h := &BlockHeader{}

  err := h.FromUntyped(fields)

  return &WrappedBlockHeader{h}, err
}

func (h *WrappedBlockHeader) FromUntyped(fields_ interface{}) error {

  fields, ok := fields_.([]byte)
  if !ok {
    return errors.New("expected byte list")
  }

  bh, err := BlockHeaderFromCBOR(fields)
  if err != nil {
    return err
  }

  h.Header = bh

  return nil
}

func (h *WrappedBlockHeader) ToUntyped() interface{} {
  var untyped interface{} = h.Header.ToCBOR()

  return untyped
}

