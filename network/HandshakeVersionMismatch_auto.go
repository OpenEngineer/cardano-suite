package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func HandshakeVersionMismatchDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func HandshakeVersionMismatchFromUntyped(fields interface{}) (*HandshakeVersionMismatch, error) {
  x := &HandshakeVersionMismatch{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *HandshakeVersionMismatch) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  ValidVersions, err := common.IntListFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for HandshakeVersionMismatch.ValidVersions: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.ValidVersions = ValidVersions
  return nil
}

func (x *HandshakeVersionMismatch) ToUntyped() interface{} {
  d := make([]interface{}, 1)
  {
    var untyped interface{} = x.ValidVersions
    d[0] = untyped
  }
  var d_ interface{} = d
  return d_
}

func HandshakeVersionMismatchFromCBOR(b []byte) (*HandshakeVersionMismatch, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return HandshakeVersionMismatchFromUntyped(d_)
}

func (x *HandshakeVersionMismatch) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
