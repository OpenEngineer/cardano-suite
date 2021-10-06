package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func BlockFetchRequestRangeDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func BlockFetchRequestRangeFromUntyped(fields interface{}) (*BlockFetchRequestRange, error) {
  x := &BlockFetchRequestRange{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *BlockFetchRequestRange) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Lower, err := ledger.PointFromUntyped(fields[0])
  if err != nil {
    return err
  }
  x.Lower = Lower
  Upper, err := ledger.PointFromUntyped(fields[1])
  if err != nil {
    return err
  }
  x.Upper = Upper
  return nil
}

func (x *BlockFetchRequestRange) ToUntyped() interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Lower.ToUntyped()
    d[0] = untyped
  }
  {
    var untyped interface{} = x.Upper.ToUntyped()
    d[1] = untyped
  }
  var d_ interface{} = d
  return d_
}

func BlockFetchRequestRangeFromCBOR(b []byte) (*BlockFetchRequestRange, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockFetchRequestRangeFromUntyped(d_)
}

func (x *BlockFetchRequestRange) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
