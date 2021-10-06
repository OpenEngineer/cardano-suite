package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func ChainSyncRollForwardDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func ChainSyncRollForwardFromUntyped(fields interface{}) (*ChainSyncRollForward, error) {
  x := &ChainSyncRollForward{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *ChainSyncRollForward) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Header, err := WrappedBlockHeaderFromUntyped(fields[0])
  if err != nil {
    return err
  }
  x.Header = Header
  Tip, err := ledger.TipFromUntyped(fields[1])
  if err != nil {
    return err
  }
  x.Tip = Tip
  return nil
}

func (x *ChainSyncRollForward) ToUntyped() interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Header.ToUntyped()
    d[0] = untyped
  }
  {
    var untyped interface{} = x.Tip.ToUntyped()
    d[1] = untyped
  }
  var d_ interface{} = d
  return d_
}

func ChainSyncRollForwardFromCBOR(b []byte) (*ChainSyncRollForward, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncRollForwardFromUntyped(d_)
}

func (x *ChainSyncRollForward) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
