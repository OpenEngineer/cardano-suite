package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func HandshakeRefuseFromUntyped(fields []interface{}) (HandshakeRefuse, error) {
  t, ok := fields[0].(uint64)
  if !ok {
    return nil, errors.New("first field entry isn't an int type but " + reflect.TypeOf(fields[0]).String())
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
      return nil, errors.New("type int out of range")
  }
}

func HandshakeRefuseToUntyped(x_ HandshakeRefuse) []interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *HandshakeVersionMismatch:
      d = append(d, int(0))
      d = append(d, x.ToUntyped()...)
    case *HandshakeDecodeError:
      d = append(d, int(1))
      d = append(d, x.ToUntyped()...)
    case *HandshakeRefused:
      d = append(d, int(2))
      d = append(d, x.ToUntyped()...)
  }
  return d
}

func HandshakeRefuseFromCBOR(b []byte) (HandshakeRefuse, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeRefuseFromUntyped(d)
}

func HandshakeRefuseToCBOR(x HandshakeRefuse) []byte {
  d := HandshakeRefuseToUntyped(x)
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
