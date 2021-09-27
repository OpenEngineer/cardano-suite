package cardano

import (
  "errors"
  cbor "https://github.com/fxamacker/cbor/v2"
)

func HandshakeAcceptVersionFromUntyped(fields []interface{}) (*HandshakeAcceptVersion, error) {
  x := &HandshakeAcceptVersion{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
}

func (x *HandshakeAcceptVersion) FromUntyped(fields []interface{}) error {
  Version, ok := fields[0].(int)
  if !ok {
    return errors.New("unexpected field type for int.Version")
  }
  x.Version = Version
  Param, ok := fields[1].(int)
  if !ok {
    return errors.New("unexpected field type for int.Param")
  }
  x.Param = Param
}

func (x *HandshakeAcceptVersion) ToUntyped() []interface{} {
  d := make([]interface{}, 2)
  d[0] = x.Version.(int)
  return d
  d[1] = x.Param.(int)
  return d
}

func HandshakeAcceptVersionFromCBOR(b []byte) (*HandshakeAcceptVersion, error) {
  d := make([]interface{}, 0)
  err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeAcceptVersionFromUntyped(d)
}

func (x *HandshakeAcceptVersion) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d, cbor.CanonicalEncOptions())
  if err != nil {
    return err
  }
  return b
}
