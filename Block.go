package cardano

import (
)

// TODO: make this an interface
type Block struct {
  Untyped interface{} // TODO: reverse engineer this, differs for every era, so Block will probably be an interface
}

type ByronBlock struct {
}

type ShelleyBlock struct {
}

type WrappedBlock struct {
  Block *Block
}

func WrappedBlockFromUntyped(fields interface{}) (*WrappedBlock, error) {
  wb := &WrappedBlock{}

  err := wb.FromUntyped(fields)

  return wb, err
}

func (b *WrappedBlock) FromUntyped(x_ interface{}) error {
  x, err := unwrapCBOR(x_)
  if err != nil {
    return err
  }

  bb := &Block{x}

  b.Block = bb

  return nil
}

func (b *WrappedBlock) ToUntyped() interface{} {
  untyped, err := wrapCBOR(b.Block.Untyped)
  if err != nil {
    panic(err)
  }

  return untyped
}
