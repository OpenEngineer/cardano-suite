package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func BlockFetchBatchDoneDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func BlockFetchBatchDoneFromUntyped(fields interface{}) (*BlockFetchBatchDone, error) {
  x := &BlockFetchBatchDone{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *BlockFetchBatchDone) FromUntyped(fields_ interface{}) error {
  return nil
}

func (x *BlockFetchBatchDone) ToUntyped() interface{} {
  d := make([]interface{}, 0)
  var d_ interface{} = d
  return d_
}

func BlockFetchBatchDoneFromCBOR(b []byte) (*BlockFetchBatchDone, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockFetchBatchDoneFromUntyped(d_)
}

func (x *BlockFetchBatchDone) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
