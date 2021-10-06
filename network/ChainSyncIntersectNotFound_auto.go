package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func ChainSyncIntersectNotFoundDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func ChainSyncIntersectNotFoundFromUntyped(fields interface{}) (*ChainSyncIntersectNotFound, error) {
  x := &ChainSyncIntersectNotFound{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *ChainSyncIntersectNotFound) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Tip, err := ledger.TipFromUntyped(fields[0])
  if err != nil {
    return err
  }
  x.Tip = Tip
  return nil
}

func (x *ChainSyncIntersectNotFound) ToUntyped() interface{} {
  d := make([]interface{}, 1)
  {
    var untyped interface{} = x.Tip.ToUntyped()
    d[0] = untyped
  }
  var d_ interface{} = d
  return d_
}

func ChainSyncIntersectNotFoundFromCBOR(b []byte) (*ChainSyncIntersectNotFound, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncIntersectNotFoundFromUntyped(d_)
}

func (x *ChainSyncIntersectNotFound) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
