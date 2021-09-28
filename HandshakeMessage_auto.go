package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func HandshakeMessageFromUntyped(fields []interface{}) (HandshakeMessage, error) {
  t, ok := fields[0].(uint64)
  if !ok {
    return nil, errors.New("first field entry isn't an int type but " + reflect.TypeOf(fields[0]).String())
  }
  args := fields[1:]
  switch t {
    case 0:
      return HandshakeProposeVersionsFromUntyped(args)
    case 1:
      return HandshakeAcceptVersionFromUntyped(args)
    case 2:
      if len(args) != 1 {
        return nil, errors.New("expected a nested list")
      }
      nestedArgs, ok := args[0].([]interface{})
      if !ok {
        return nil, errors.New("expected a nested list")
      }
      return HandshakeRefuseFromUntyped(nestedArgs)
    default:
      return nil, errors.New("type int out of range")
  }
}

func HandshakeMessageToUntyped(x_ HandshakeMessage) []interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *HandshakeProposeVersions:
      d = append(d, int(0))
      d = append(d, x.ToUntyped()...)
    case *HandshakeAcceptVersion:
      d = append(d, int(1))
      d = append(d, x.ToUntyped()...)
    case HandshakeRefuse:
      d = append(d, int(2))
      d = append(d, HandshakeRefuseToUntyped(x))
  }
  return d
}

func HandshakeMessageFromCBOR(b []byte) (HandshakeMessage, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeMessageFromUntyped(d)
}

func HandshakeMessageToCBOR(x HandshakeMessage) []byte {
  d := HandshakeMessageToUntyped(x)
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
