package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func HandshakeRefusedDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func HandshakeRefusedFromUntyped(fields interface{}) (*HandshakeRefused, error) {
  x := &HandshakeRefused{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *HandshakeRefused) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Version, err := common.IntFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for HandshakeRefused.Version: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.Version = Version
  Reason, err := common.StringFromUntyped(fields[1])
  if err != nil {
    return errors.New("unexpected field type for HandshakeRefused.Reason: " + reflect.TypeOf(fields[1]).String() + " " + err.Error())
  }
  x.Reason = Reason
  return nil
}

func (x *HandshakeRefused) ToUntyped() interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Version
    d[0] = untyped
  }
  {
    var untyped interface{} = x.Reason
    d[1] = untyped
  }
  var d_ interface{} = d
  return d_
}

func HandshakeRefusedFromCBOR(b []byte) (*HandshakeRefused, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return HandshakeRefusedFromUntyped(d_)
}

func (x *HandshakeRefused) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
