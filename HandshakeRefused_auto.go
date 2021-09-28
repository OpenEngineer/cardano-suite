package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func HandshakeRefusedFromUntyped(fields []interface{}) (*HandshakeRefused, error) {
  x := &HandshakeRefused{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *HandshakeRefused) FromUntyped(fields []interface{}) error {
  Version, err := IntFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for HandshakeRefused.Version: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.Version = Version
  Error, err := StringFromUntyped(fields[1])
  if err != nil {
    return errors.New("unexpected field type for HandshakeRefused.Error: " + reflect.TypeOf(fields[1]).String() + " " + err.Error())
  }
  x.Error = Error
  return nil
}

func (x *HandshakeRefused) ToUntyped() []interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Version
    d[0] = untyped
  }
  {
    var untyped interface{} = x.Error
    d[1] = untyped
  }
  return d
}

func HandshakeRefusedFromCBOR(b []byte) (*HandshakeRefused, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeRefusedFromUntyped(d)
}

func (x *HandshakeRefused) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
