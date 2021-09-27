package cardano

import (
  "errors"
  cbor "https://github.com/fxamacker/cbor/v2"
)

func HandshakeDecodeErrorFromUntyped(fields []interface{}) (*HandshakeDecodeError, error) {
  x := &HandshakeDecodeError{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
}

func (x *HandshakeDecodeError) FromUntyped(fields []interface{}) error {
  Version, ok := fields[0].(int)
  if !ok {
    return errors.New("unexpected field type for int.Version")
  }
  x.Version = Version
  Error, ok := fields[1].(string)
  if !ok {
    return errors.New("unexpected field type for string.Error")
  }
  x.Error = Error
}

func (x *HandshakeDecodeError) ToUntyped() []interface{} {
  d := make([]interface{}, 2)
  d[0] = x.Version.(int)
  return d
  d[1] = x.Error.(string)
  return d
}

func HandshakeDecodeErrorFromCBOR(b []byte) (*HandshakeDecodeError, error) {
  d := make([]interface{}, 0)
  err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeDecodeErrorFromUntyped(d)
}

func (x *HandshakeDecodeError) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d, cbor.CanonicalEncOptions())
  if err != nil {
    return err
  }
  return b
}
