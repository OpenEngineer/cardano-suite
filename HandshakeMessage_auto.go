package cardano

import (
  "errors"
  cbor "https://github.com/fxamacker/cbor/v2"
)

func HandshakeMessageFromUntyped(fields []interface{}) (HandshakeMessage, error) {
  t, ok := fields[0].(byte)
  if !ok {
    return nil, errors.New("first field entry isn't a type byte")
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
      return nil, errors.New("type byte out of range")
  }
}

func HandshakeMessageToUntyped(x_ HandshakeMessage) []interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *HandshakeProposeVersions:
      d = append(d, byte(0))
      d = append(d, x.ToUntyped()...)
    case *HandshakeAcceptVersion:
      d = append(d, byte(1))
      d = append(d, x.ToUntyped()...)
    case HandshakeRefuse:
      d = append(d, byte(2))
      d = append(d, HandshakeRefuseToUntyped(x))
  }
  return d
}

func HandshakeMessageFromCBOR(b []byte) (HandshakeMessage, error) {
  d := make([]interface{}, 0)
  err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  return HandshakeMessageFromUntyped(d)
}

func HandshakeMessageToCBOR(x HandshakeMessage) []byte {
  d := HandshakeMessageToUntyped(x)
  b, err := cbor.Marshal(d, cbor.CanonicalEncOptions())
  if err != nil {
    return err
  }
  return b
}
