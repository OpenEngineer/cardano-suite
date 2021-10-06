package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func BlockFetchBlockDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func BlockFetchBlockFromUntyped(fields interface{}) (*BlockFetchBlock, error) {
  x := &BlockFetchBlock{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *BlockFetchBlock) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Block, err := WrappedBlockFromUntyped(fields[0])
  if err != nil {
    return err
  }
  x.Block = Block
  return nil
}

func (x *BlockFetchBlock) ToUntyped() interface{} {
  d := make([]interface{}, 1)
  {
    var untyped interface{} = x.Block.ToUntyped()
    d[0] = untyped
  }
  var d_ interface{} = d
  return d_
}

func BlockFetchBlockFromCBOR(b []byte) (*BlockFetchBlock, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockFetchBlockFromUntyped(d_)
}

func (x *BlockFetchBlock) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
