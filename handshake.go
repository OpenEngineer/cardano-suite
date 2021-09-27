package cardano

type HandshakeMessage interface {
  ToUntyped() []interface{}
}

type HandshakeProposeVersions struct {
  Versions map[int]int // map acts as a set
}

type HandshakeAcceptVersion struct {
  Version int
  Param   int
}

type HandshakeRefuse interface {
  HandshakeMessage
}

type HandshakeVersionMismatch struct {
  ValidVersions []int
}

type HandshakeDecodeError struct {
  Version int
  Error   string
}

type HandshakeRefused struct {
  Version int
  Error   string
}

//go:generate ./cbor-type HandshakeMessage *HandshakeProposeVersions *HandshakeAcceptVersion HandshakeRefuse
//go:generate ./cbor-type *HandshakeProposeVersions "Versions map[int]int"
//go:generate ./cbor-type *HandshakeAcceptVersion   "Version int, Param int"
//go:generate ./cbor-type HandshakeRefuse *HandshakeVersionMismatch *HandshakeDecodeError *HandshakeRefused
//go:generate ./cbor-type *HandshakeVersionMismatch "ValidVersions []int"
//go:generate ./cbor-type *HandshakeDecodeError     "Version int, Error string"
//go:generate ./cbor-type *HandshakeRefused         "Version int, Error string"
