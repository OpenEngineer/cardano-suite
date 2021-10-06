package ledger

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
)

func BlockHeaderDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String())}

func BlockHeaderFromUntyped(fields interface{}) (*BlockHeader, error) {
  x := &BlockHeader{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *BlockHeader) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  HeaderHash, err := common.HashFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for BlockHeader.HeaderHash: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.HeaderHash = HeaderHash
  PrevBlock, err := ChainHashFromUntyped(fields[1])
  if err != nil {
    return errors.New("unexpected field type for BlockHeader.PrevBlock: " + reflect.TypeOf(fields[1]).String() + " " + err.Error())
  }
  x.PrevBlock = PrevBlock
  SlotId, err := common.Uint64FromUntyped(fields[2])
  if err != nil {
    return errors.New("unexpected field type for BlockHeader.SlotId: " + reflect.TypeOf(fields[2]).String() + " " + err.Error())
  }
  x.SlotId = SlotId
  BlockId, err := common.Uint64FromUntyped(fields[3])
  if err != nil {
    return errors.New("unexpected field type for BlockHeader.BlockId: " + reflect.TypeOf(fields[3]).String() + " " + err.Error())
  }
  x.BlockId = BlockId
  BodyHash, err := common.HashFromUntyped(fields[4])
  if err != nil {
    return errors.New("unexpected field type for BlockHeader.BodyHash: " + reflect.TypeOf(fields[4]).String() + " " + err.Error())
  }
  x.BodyHash = BodyHash
  return nil
}

func (x *BlockHeader) ToUntyped() interface{} {
  d := make([]interface{}, 5)
  {
    var untyped interface{} = x.HeaderHash.ToUntyped()
    d[0] = untyped
  }
  {
    var untyped interface{} = x.PrevBlock.ToUntyped()
    d[1] = untyped
  }
  {
    var untyped interface{} = x.SlotId
    d[2] = untyped
  }
  {
    var untyped interface{} = x.BlockId
    d[3] = untyped
  }
  {
    var untyped interface{} = x.BodyHash.ToUntyped()
    d[4] = untyped
  }
  var d_ interface{} = d
  return d_
}

func BlockHeaderFromCBOR(b []byte) (*BlockHeader, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockHeaderFromUntyped(d_)
}

func (x *BlockHeader) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
