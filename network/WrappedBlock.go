package network

import (
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

type WrappedBlock struct {
  // Block should be an interface
  Block *ledger.Block
}

func WrappedBlockFromUntyped(fields interface{}) (*WrappedBlock, error) {
  wb := &WrappedBlock{}

  err := wb.FromUntyped(fields)

  return wb, err
}

func (b *WrappedBlock) FromUntyped(x_ interface{}) error {
  x, err := common.UnwrapCBOR(x_)
  if err != nil {
    return err
  }

  bb := &ledger.Block{x}

  b.Block = bb

  return nil
}

func (b *WrappedBlock) ToUntyped() interface{} {
  untyped, err := common.WrapCBOR(b.Block.Untyped)
  if err != nil {
    panic(err)
  }

  return untyped
}
