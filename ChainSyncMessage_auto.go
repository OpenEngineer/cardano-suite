package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func ChainSyncMessageDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func ChainSyncMessageFromUntyped(fields_ interface{}) (ChainSyncMessage, error) {
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
      return ChainSyncRequestNextFromUntyped(args)
    case 1:
      return ChainSyncAwaitReplyFromUntyped(args)
    case 2:
      return ChainSyncRollForwardFromUntyped(args)
    case 3:
      return ChainSyncRollBackwardFromUntyped(args)
    case 4:
      return ChainSyncFindIntersectFromUntyped(args)
    case 5:
      return ChainSyncIntersectFoundFromUntyped(args)
    case 6:
      return ChainSyncIntersectNotFoundFromUntyped(args)
    case 7:
      return ChainSyncDoneFromUntyped(args)
    default:
      return nil, errors.New("type int out of range")
  }
}

func ChainSyncMessageToUntyped(x_ ChainSyncMessage) interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *ChainSyncRequestNext:
      d = append(d, int(0))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *ChainSyncAwaitReply:
      d = append(d, int(1))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *ChainSyncRollForward:
      d = append(d, int(2))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *ChainSyncRollBackward:
      d = append(d, int(3))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *ChainSyncFindIntersect:
      d = append(d, int(4))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *ChainSyncIntersectFound:
      d = append(d, int(5))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *ChainSyncIntersectNotFound:
      d = append(d, int(6))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *ChainSyncDone:
      d = append(d, int(7))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
  }
  var d_ interface{} = d
  return d_
}

func ChainSyncMessageFromCBOR(b []byte) (ChainSyncMessage, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncMessageFromUntyped(d_)
}

func ChainSyncMessageToCBOR(x ChainSyncMessage) []byte {
  d := ChainSyncMessageToUntyped(x)
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
