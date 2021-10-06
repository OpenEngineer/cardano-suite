package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func BlockFetchNoBlocksDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func BlockFetchNoBlocksFromUntyped(fields interface{}) (*BlockFetchNoBlocks, error) {
  x := &BlockFetchNoBlocks{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *BlockFetchNoBlocks) FromUntyped(fields_ interface{}) error {
  return nil
}

func (x *BlockFetchNoBlocks) ToUntyped() interface{} {
  d := make([]interface{}, 0)
  var d_ interface{} = d
  return d_
}

func BlockFetchNoBlocksFromCBOR(b []byte) (*BlockFetchNoBlocks, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockFetchNoBlocksFromUntyped(d_)
}

func (x *BlockFetchNoBlocks) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
