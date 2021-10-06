package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func BlockFetchStartBatchDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func BlockFetchStartBatchFromUntyped(fields interface{}) (*BlockFetchStartBatch, error) {
  x := &BlockFetchStartBatch{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *BlockFetchStartBatch) FromUntyped(fields_ interface{}) error {
  return nil
}

func (x *BlockFetchStartBatch) ToUntyped() interface{} {
  d := make([]interface{}, 0)
  var d_ interface{} = d
  return d_
}

func BlockFetchStartBatchFromCBOR(b []byte) (*BlockFetchStartBatch, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockFetchStartBatchFromUntyped(d_)
}

func (x *BlockFetchStartBatch) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
