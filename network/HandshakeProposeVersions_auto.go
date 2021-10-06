package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func HandshakeProposeVersionsDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func HandshakeProposeVersionsFromUntyped(fields interface{}) (*HandshakeProposeVersions, error) {
  x := &HandshakeProposeVersions{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *HandshakeProposeVersions) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Versions, err := common.IntToInterfMapFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for HandshakeProposeVersions.Versions: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.Versions = Versions
  return nil
}

func (x *HandshakeProposeVersions) ToUntyped() interface{} {
  d := make([]interface{}, 1)
  {
    var untyped interface{} = x.Versions
    d[0] = untyped
  }
  var d_ interface{} = d
  return d_
}

func HandshakeProposeVersionsFromCBOR(b []byte) (*HandshakeProposeVersions, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return HandshakeProposeVersionsFromUntyped(d_)
}

func (x *HandshakeProposeVersions) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
