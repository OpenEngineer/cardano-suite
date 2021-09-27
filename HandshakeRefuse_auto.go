package cardano

import (
  "errors"
  cbor "https://github.com/fxamacker/cbor/v2"
)

func HandshakeRefuseFromUntyped(fields []interface{}) (HandshakeRefuse, error) {
  t, ok := fields[0].(byte)
  if !ok {
    return nil, errors.New("first field entry isn't a type byte")
  }
  args := fields[1:]
  switch t {
    case 0:
      return HandshakeVersionMismatchFromUntyped(args)
    case 1:
      return HandshakeDecodeErrorFromUntyped(args)
    case 2:
      return HandshakeRefusedFromUntyped(args)
    default:
      return nil, errors.New("type byte out of range")
  }
}

func HandshakeRefuseToUntyped(x_ HandshakeRefuse) []interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *HandshakeVersionMismatch:
      d = append(d, byte(0))
      d = append(d, x.ToUntyped()...)
    case *HandshakeDecodeError:
      d = append(d, byte(1))
      d = append(d, x.ToUntyped()...)
    case *HandshakeRefused:
      d = append(d, byte(2))
      d = append(d, x.ToUntyped()...)
  }
  return d
}

func HandshakeRefuseFromCBOR(b []byte) (HandshakeRefuse, error) {
  d := make([]interface{}, 0)
  err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeRefuseFromUntyped(d)
}

func HandshakeRefuseToCBOR(x HandshakeRefuse) []byte {
  d := HandshakeRefuseToUntyped(x)
  b, err := cbor.Marshal(d, cbor.CanonicalEncOptions())
  if err != nil {
    return err
  }
  return b
}
