package ledger

import (
	"errors"

	"github.com/christianschmitz/cardano-suite/common"
)

// compact representation of BlockHeader for ChainSync protocol
type Point struct {
	IsGenesis  bool `emptyiftrue`
	SlotId     uint64
	HeaderHash common.Hash
}

func PointFromUntyped(fields interface{}) (*Point, error) {
	p := &Point{}
	err := p.FromUntyped(fields)
	return p, err
}

func (p *Point) FromUntyped(fields_ interface{}) error {
	fields, err := common.InterfListFromUntyped(fields_)
	if err != nil {
		return err
	}

	switch len(fields) {
	case 0:
		p.IsGenesis = true
	case 2:
		sId, err := common.Uint64FromUntyped(fields[0])
		if err != nil {
			return err
		}

		p.SlotId = sId

		h, err := common.HashFromUntyped(fields[1])
		if err != nil {
			return err
		}

		p.HeaderHash = h
	default:
		return errors.New("expected 0 or 2 fields")
	}

	return nil
}

func (p *Point) ToUntyped() interface{} {
	res := make([]interface{}, 0)

	if !p.IsGenesis {
		res = append(res, p.SlotId, p.HeaderHash.ToUntyped())
	}

	var untyped interface{} = res

	return untyped
}
