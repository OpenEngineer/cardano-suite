package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func ChainSyncAwaitReplyDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func ChainSyncAwaitReplyFromUntyped(fields interface{}) (*ChainSyncAwaitReply, error) {
  x := &ChainSyncAwaitReply{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *ChainSyncAwaitReply) FromUntyped(fields_ interface{}) error {
  return nil
}

func (x *ChainSyncAwaitReply) ToUntyped() interface{} {
  d := make([]interface{}, 0)
  var d_ interface{} = d
  return d_
}

func ChainSyncAwaitReplyFromCBOR(b []byte) (*ChainSyncAwaitReply, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncAwaitReplyFromUntyped(d_)
}

func (x *ChainSyncAwaitReply) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
