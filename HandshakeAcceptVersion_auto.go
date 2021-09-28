package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func HandshakeAcceptVersionFromUntyped(fields []interface{}) (*HandshakeAcceptVersion, error) {
  x := &HandshakeAcceptVersion{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *HandshakeAcceptVersion) FromUntyped(fields []interface{}) error {
  Version, err := IntFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for HandshakeAcceptVersion.Version: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.Version = Version
  Param, err := InterfFromUntyped(fields[1])
  if err != nil {
    return errors.New("unexpected field type for HandshakeAcceptVersion.Param: " + reflect.TypeOf(fields[1]).String() + " " + err.Error())
  }
  x.Param = Param
  return nil
}

func (x *HandshakeAcceptVersion) ToUntyped() []interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Version
    d[0] = untyped
  }
  {
    var untyped interface{} = x.Param
    d[1] = untyped
  }
  return d
}

func HandshakeAcceptVersionFromCBOR(b []byte) (*HandshakeAcceptVersion, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeAcceptVersionFromUntyped(d)
}

func (x *HandshakeAcceptVersion) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
