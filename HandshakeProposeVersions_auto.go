package cardano

import (
  "errors"
  cbor "https://github.com/fxamacker/cbor/v2"
)

func HandshakeProposeVersionsFromUntyped(fields []interface{}) (*HandshakeProposeVersions, error) {
  x := &HandshakeProposeVersions{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
}

func (x *HandshakeProposeVersions) FromUntyped(fields []interface{}) error {
  Versions, ok := fields[0].(map[int]int)
  if !ok {
    return errors.New("unexpected field type for map[int]int.Versions")
  }
  x.Versions = Versions
}

func (x *HandshakeProposeVersions) ToUntyped() []interface{} {
  d := make([]interface{}, 1)
  d[0] = x.Versions.(map[int]int)
  return d
}

func HandshakeProposeVersionsFromCBOR(b []byte) (*HandshakeProposeVersions, error) {
  d := make([]interface{}, 0)
  err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeProposeVersionsFromUntyped(d)
}

func (x *HandshakeProposeVersions) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d, cbor.CanonicalEncOptions())
  if err != nil {
    return err
  }
  return b
}
