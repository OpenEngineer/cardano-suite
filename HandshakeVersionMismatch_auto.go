package cardano

import (
  "errors"
  cbor "https://github.com/fxamacker/cbor/v2"
)

func HandshakeVersionMismatchFromUntyped(fields []interface{}) (*HandshakeVersionMismatch, error) {
  x := &HandshakeVersionMismatch{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
}

func (x *HandshakeVersionMismatch) FromUntyped(fields []interface{}) error {
  ValidVersions, ok := fields[0].([]int)
  if !ok {
    return errors.New("unexpected field type for []int.ValidVersions")
  }
  x.ValidVersions = ValidVersions
}

func (x *HandshakeVersionMismatch) ToUntyped() []interface{} {
  d := make([]interface{}, 1)
  d[0] = x.ValidVersions.([]int)
  return d
}

func HandshakeVersionMismatchFromCBOR(b []byte) (*HandshakeVersionMismatch, error) {
  d := make([]interface{}, 0)
  err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeVersionMismatchFromUntyped(d)
}

func (x *HandshakeVersionMismatch) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d, cbor.CanonicalEncOptions())
  if err != nil {
    return err
  }
  return b
}
