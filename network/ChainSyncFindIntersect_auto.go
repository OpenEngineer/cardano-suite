package network

import (
  "errors"
  "reflect"
  cbor "github.com/fxamacker/cbor/v2"
  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func ChainSyncFindIntersectDummyImportUsage() error {return errors.New(reflect.TypeOf(common.Hash{}).String() + reflect.TypeOf(ledger.Block{}).String())}

func ChainSyncFindIntersectFromUntyped(fields interface{}) (*ChainSyncFindIntersect, error) {
  x := &ChainSyncFindIntersect{}
  if err := x.FromUntyped(fields); err != nil {
    return nil, err
  }
  return x, nil
}

func (x *ChainSyncFindIntersect) FromUntyped(fields_ interface{}) error {
  fields, err := common.InterfListFromUntyped(fields_)
  if err != nil {
    return err
  }
  PointsList, ok := fields[0].([]interface{})
  if !ok {
    return errors.New("expected []*ledger.Point for ChainSyncFindIntersect.Points")
  }
  Points := make([]*ledger.Point, len(PointsList))
  for i, item := range PointsList {
    Points[i], err = ledger.PointFromUntyped(item)
    if err != nil {
      return err
    }
  }
  x.Points = Points
  return nil
}

func (x *ChainSyncFindIntersect) ToUntyped() interface{} {
  d := make([]interface{}, 1)
  {
    lst := make([]interface{}, len(x.Points))
    for i, item := range x.Points {
      lst[i] = item.ToUntyped()
    }
    var untyped interface{} = lst
    d[0] = untyped
  }
  var d_ interface{} = d
  return d_
}

func ChainSyncFindIntersectFromCBOR(b []byte) (*ChainSyncFindIntersect, error) {
  d := make([]interface{}, 0)
  if err := cbor.Unmarshal(b, &d); err != nil {
    return nil, err
  }
  var d_ interface{} = d
  return ChainSyncFindIntersectFromUntyped(d_)
}

func (x *ChainSyncFindIntersect) ToCBOR() []byte {
  d := x.ToUntyped()
  b, err := cbor.Marshal(d)
  if err != nil {
    panic(err)
  }
  return b
}
