package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func ChainSyncDoneDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func ChainSyncDoneFromUntyped(fields interface{}) (*ChainSyncDone, error) {
  x := &ChainSyncDone{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *ChainSyncDone) FromUntyped(fields_ interface{}) error {
  return nil
}

func (x *ChainSyncDone) ToUntyped() interface{} {
  d := make([]interface{}, 0)
  var d_ interface{} = d
  return d_
}

func ChainSyncDoneFromCBOR(b []byte) (*ChainSyncDone, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncDoneFromUntyped(d_)
}

func (x *ChainSyncDone) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
