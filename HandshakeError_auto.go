package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func HandshakeErrorDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func HandshakeErrorFromUntyped(fields_ interface{}) (HandshakeError, error) {
  fields, err := InterfListFromUntyped(fields_)
  if err != nil {
    return nil, err
  }
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

func HandshakeErrorToUntyped(x_ HandshakeError) interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *HandshakeVersionMismatch:
      d = append(d, int(0))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *HandshakeDecodeError:
      d = append(d, int(1))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *HandshakeRefused:
      d = append(d, int(2))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
  }
  var d_ interface{} = d
  return d_
}

func HandshakeErrorFromCBOR(b []byte) (HandshakeError, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return HandshakeErrorFromUntyped(d_)
}

func HandshakeErrorToCBOR(x HandshakeError) []byte {
  d := HandshakeErrorToUntyped(x)
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
