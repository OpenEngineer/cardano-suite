package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func BlockFetchClientDoneDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func BlockFetchClientDoneFromUntyped(fields interface{}) (*BlockFetchClientDone, error) {
  x := &BlockFetchClientDone{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *BlockFetchClientDone) FromUntyped(fields_ interface{}) error {
  return nil
}

func (x *BlockFetchClientDone) ToUntyped() interface{} {
  d := make([]interface{}, 0)
  var d_ interface{} = d
  return d_
}

func BlockFetchClientDoneFromCBOR(b []byte) (*BlockFetchClientDone, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return BlockFetchClientDoneFromUntyped(d_)
}

func (x *BlockFetchClientDone) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
