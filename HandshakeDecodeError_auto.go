package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func HandshakeDecodeErrorFromUntyped(fields []interface{}) (*HandshakeDecodeError, error) {
  x := &HandshakeDecodeError{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *HandshakeDecodeError) FromUntyped(fields []interface{}) error {
  Version, err := IntFromUntyped(fields[0])
  if err != nil {
    return errors.New("unexpected field type for HandshakeDecodeError.Version: " + reflect.TypeOf(fields[0]).String() + " " + err.Error())
  }
  x.Version = Version
  Error, err := StringFromUntyped(fields[1])
  if err != nil {
    return errors.New("unexpected field type for HandshakeDecodeError.Error: " + reflect.TypeOf(fields[1]).String() + " " + err.Error())
  }
  x.Error = Error
  return nil
}

func (x *HandshakeDecodeError) ToUntyped() []interface{} {
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

func HandshakeDecodeErrorFromCBOR(b []byte) (*HandshakeDecodeError, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeDecodeErrorFromUntyped(d)
}

func (x *HandshakeDecodeError) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
