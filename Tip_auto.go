package cardano

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
)

func TipDummyImportUsage() error {return errors.New(reflect.TypeOf("").String())}

func TipFromUntyped(fields interface{}) (*Tip, error) {
  x := &Tip{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *Tip) FromUntyped(fields_ interface{}) error {
  fields, err := InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  Point, err := PointFromUntyped(fields[0])
  if err != nil {
    return err
  }
  x.Point = Point
  BlockId, err := Uint64FromUntyped(fields[1])
  if err != nil {
    return errors.New("unexpected field type for Tip.BlockId: " + reflect.TypeOf(fields[1]).String() + " " + err.Error())
  }
  x.BlockId = BlockId
  return nil
}

func (x *Tip) ToUntyped() interface{} {
  d := make([]interface{}, 2)
  {
    var untyped interface{} = x.Point.ToUntyped()
    d[0] = untyped
  }
  {
    var untyped interface{} = x.BlockId
    d[1] = untyped
  }
  var d_ interface{} = d
  return d_
}

func TipFromCBOR(b []byte) (*Tip, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return TipFromUntyped(d_)
}

func (x *Tip) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
