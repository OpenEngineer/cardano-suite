package cardano

import (
  "errors"
  cbor "https://github.com/fxamacker/cbor/v2"
)

func HandshakeRefusedFromUntyped(fields []interface{}) (*HandshakeRefused, error) {
  x := &HandshakeRefused{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
}

func (x *HandshakeRefused) FromUntyped(fields []interface{}) error {
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

func (x *HandshakeRefused) ToUntyped() []interface{} {
  d := make([]interface{}, 2)
  d[0] = x.Version.(int)
  return d
  d[1] = x.Error.(string)
  return d
}

func HandshakeRefusedFromCBOR(b []byte) (*HandshakeRefused, error) {
  d := make([]interface{}, 0)
  err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeRefusedFromUntyped(d)
}

func (x *HandshakeRefused) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d, cbor.CanonicalEncOptions())
  if err != nil {
    return err
  }
  return b
}
