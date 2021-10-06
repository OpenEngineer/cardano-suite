package ledger

import (
)

// TODO: make this an interface
type Block struct {
  Untyped interface{} // TODO: reverse engineer this, differs for every era, so Block will probably be an interface
}

type ByronBlock struct {
}

type ShelleyBlock struct {
}
