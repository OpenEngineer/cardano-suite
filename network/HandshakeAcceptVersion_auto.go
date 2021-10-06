package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func HandshakeAcceptVersionDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func HandshakeAcceptVersionFromUntyped(fields interface{}) (*HandshakeAcceptVersion, error) {
  x := &HandshakeAcceptVersion{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *HandshakeAcceptVersion) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Version, err := common.IntFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for HandshakeAcceptVersion.Version: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.Version = Version
  Params, err := common.InterfFromUntyped(fields[1])
  if err != nil {
    return errors.New("unexpected field type for HandshakeAcceptVersion.Params: " + reflect.TypeOf(fields[1]).String() + " " + err.Error())
  }
  x.Params = Params
  return nil
}

func (x *HandshakeAcceptVersion) ToUntyped() interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Version
    d[0] = untyped
  }
  {
    var untyped interface{} = x.Params
    d[1] = untyped
  }
  var d_ interface{} = d
  return d_
}

func HandshakeAcceptVersionFromCBOR(b []byte) (*HandshakeAcceptVersion, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return HandshakeAcceptVersionFromUntyped(d_)
}

func (x *HandshakeAcceptVersion) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
