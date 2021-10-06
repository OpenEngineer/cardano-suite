package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func HandshakeMessageDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func HandshakeMessageFromUntyped(fields_ interface{}) (HandshakeMessage, error) {
  fields, err := common.InterfListFromUntyped(fields_)
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
      return HandshakeProposeVersionsFromUntyped(args)
    case 1:
      return HandshakeAcceptVersionFromUntyped(args)
    case 2:
      if len(args) != 1 {
        return nil, errors.New("expected a nested list")
      }
      return HandshakeErrorFromUntyped(args[0])
    default:
      return nil, errors.New("type int out of range")
  }
}

func HandshakeMessageToUntyped(x_ HandshakeMessage) interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *HandshakeProposeVersions:
      d = append(d, int(0))
      tmpLst, err := common.InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *HandshakeAcceptVersion:
      d = append(d, int(1))
      tmpLst, err := common.InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case HandshakeError:
      d = append(d, int(2))
      d = append(d, HandshakeErrorToUntyped(x))
  }
  var d_ interface{} = d
  return d_
}

func HandshakeMessageFromCBOR(b []byte) (HandshakeMessage, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return HandshakeMessageFromUntyped(d_)
}

func HandshakeMessageToCBOR(x HandshakeMessage) []byte {
  d := HandshakeMessageToUntyped(x)
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
