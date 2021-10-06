package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func BlockFetchMessageDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func BlockFetchMessageFromUntyped(fields_ interface{}) (BlockFetchMessage, error) {
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
      return BlockFetchRequestRangeFromUntyped(args)
    case 1:
      return BlockFetchClientDoneFromUntyped(args)
    case 2:
      return BlockFetchStartBatchFromUntyped(args)
    case 3:
      return BlockFetchNoBlocksFromUntyped(args)
    case 4:
      return BlockFetchBlockFromUntyped(args)
    case 5:
      return BlockFetchBatchDoneFromUntyped(args)
    default:
      return nil, errors.New("type int out of range")
  }
}

func BlockFetchMessageToUntyped(x_ BlockFetchMessage) interface{} {
  d := make([]interface{}, 0)
  switch x := x_.(type) {
    case *BlockFetchRequestRange:
      d = append(d, int(0))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *BlockFetchClientDone:
      d = append(d, int(1))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *BlockFetchStartBatch:
      d = append(d, int(2))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *BlockFetchNoBlocks:
      d = append(d, int(3))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *BlockFetchBlock:
      d = append(d, int(4))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
    case *BlockFetchBatchDone:
      d = append(d, int(5))
      tmpLst, err := InterfListFromUntyped(x.ToUntyped())
      if err != nil {
        panic("expected []interface{} from struct to untyped")
      }
      d = append(d, tmpLst...)
  }
  var d_ interface{} = d
  return d_
}

func BlockFetchMessageFromCBOR(b []byte) (BlockFetchMessage, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockFetchMessageFromUntyped(d_)
}

func BlockFetchMessageToCBOR(x BlockFetchMessage) []byte {
  d := BlockFetchMessageToUntyped(x)
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
