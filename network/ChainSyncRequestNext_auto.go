package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func ChainSyncRequestNextDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func ChainSyncRequestNextFromUntyped(fields interface{}) (*ChainSyncRequestNext, error) {
  x := &ChainSyncRequestNext{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *ChainSyncRequestNext) FromUntyped(fields_ interface{}) error {
  return nil
}

func (x *ChainSyncRequestNext) ToUntyped() interface{} {
  d := make([]interface{}, 0)
  var d_ interface{} = d
  return d_
}

func ChainSyncRequestNextFromCBOR(b []byte) (*ChainSyncRequestNext, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncRequestNextFromUntyped(d_)
}

func (x *ChainSyncRequestNext) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
