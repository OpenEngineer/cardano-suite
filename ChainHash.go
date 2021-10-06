package cardano 

import (
  "errors"
)

// used in BlockHeader to point to previous block 
type ChainHash struct {
  IsGenesis bool
  BlockHash Hash
}

func NewGenesisHash() ChainHash {
  return ChainHash{true, NilHash()}
}

func NewChainHash(blockHash Hash) ChainHash {
  return ChainHash{false, blockHash} // XXX: is blockHash not just headerHash of the previous block?
}

func ChainHashFromUntyped(x_ interface{}) (ChainHash, error) {
  ch := ChainHash{}

  x, err := InterfListFromUntyped(x_)
  if err != nil {
    return ch, err
  }

  switch len(x) {
  case 0:
    ch.IsGenesis = true
    ch.BlockHash = NilHash()
    return ch, nil
  case 1:
    h, err := HashFromUntyped(x[0])
    if err != nil {
      return ch, err
    }

    ch.BlockHash = h

    return ch, nil
  default:
    return ch, errors.New("invalid chainHash (MaybeGenesisHash)")
  }
}

func (h ChainHash) ToUntyped() interface{} {
  res := make([]interface{}, 0)

  if !h.IsGenesis {
    res = append(res, h.BlockHash.ToUntyped())
  }

  return res
}

func (h ChainHash) Eq(other ChainHash) bool {
  if h.IsGenesis {
    return other.IsGenesis
  } else {
    return h.BlockHash.Eq(other.BlockHash)
  }
}
