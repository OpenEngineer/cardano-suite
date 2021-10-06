package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func ChainSyncRollBackwardDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func ChainSyncRollBackwardFromUntyped(fields interface{}) (*ChainSyncRollBackward, error) {
  x := &ChainSyncRollBackward{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *ChainSyncRollBackward) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Point, err := ledger.PointFromUntyped(fields[0])
  if err != nil {
    return err
  }
  x.Point = Point
  Tip, err := ledger.TipFromUntyped(fields[1])
  if err != nil {
    return err
  }
  x.Tip = Tip
  return nil
}

func (x *ChainSyncRollBackward) ToUntyped() interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Point.ToUntyped()
    d[0] = untyped
  }
  {
    var untyped interface{} = x.Tip.ToUntyped()
    d[1] = untyped
  }
  var d_ interface{} = d
  return d_
}

func ChainSyncRollBackwardFromCBOR(b []byte) (*ChainSyncRollBackward, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncRollBackwardFromUntyped(d_)
}

func (x *ChainSyncRollBackward) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
