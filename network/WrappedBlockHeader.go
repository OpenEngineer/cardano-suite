package network

import (
  "errors"

  "github.com/christianschmitz/cardano-suite/ledger"
)

type WrappedBlockHeader struct {
  Header *ledger.BlockHeader
}

func WrappedBlockHeaderFromUntyped(fields interface{}) (*WrappedBlockHeader, error) {
  h := &ledger.BlockHeader{}

  err := h.FromUntyped(fields)

  return &WrappedBlockHeader{h}, err
}

func (h *WrappedBlockHeader) FromUntyped(fields_ interface{}) error {

  fields, ok := fields_.([]byte)
  if !ok {
    return errors.New("expected byte list")
  }

  bh, err := ledger.BlockHeaderFromCBOR(fields)
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
