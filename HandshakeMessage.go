package cardano

import (
  "errors"
)

const MAINNET_MAGIC = 764824073
const TESTNET_MAGIC = 1097911063

type HandshakeMessage interface { // all messages of a given protocol need to have the same interface type
}

type HandshakeProposeVersions struct {
  Versions map[int]interface{}
}

type HandshakeAcceptVersion struct {
  Version int
  Params   interface{}
}

type HandshakeError interface {
  HandshakeMessage
  Error() error
}

type HandshakeVersionMismatch struct {
  ValidVersions []int
}

type HandshakeDecodeError struct {
  Version int
  Reason  string
}

type HandshakeRefused struct {
  Version int
  Reason  string
}

// NOTE: the primitive type-names below are different, because the builtin cast functions are not very robust
//    see ./cbor_primitive_types.go for the possible cast functions

//go:generate ./cbor-type HandshakeMessage *HandshakeProposeVersions *HandshakeAcceptVersion HandshakeError
//go:generate ./cbor-type *HandshakeProposeVersions "Versions IntToInterfMap"
//go:generate ./cbor-type *HandshakeAcceptVersion   "Version Int, Params Interf"
//go:generate ./cbor-type HandshakeError *HandshakeVersionMismatch *HandshakeDecodeError *HandshakeRefused
//go:generate ./cbor-type *HandshakeVersionMismatch "ValidVersions IntList"
//go:generate ./cbor-type *HandshakeDecodeError     "Version Int, Reason String"
//go:generate ./cbor-type *HandshakeRefused         "Version Int, Reason String"

func NewHandshakeProposeVersions(magic int) *HandshakeProposeVersions {
  return &HandshakeProposeVersions{DefaultProtocolVersions(magic)}
}

func (m *HandshakeProposeVersions) Intersect(version map[int]interface{}) (int, interface{}, error) {
  vIntersect := -1
  var pIntersect interface{} = nil

  for v, p := range m.Versions {
    if _, ok := version[v]; ok {
      if v > vIntersect {
        vIntersect = v
        pIntersect = p
      }
    }
  }

  if vIntersect == -1 {
    return 0, nil, errors.New("no version intersection found")
  }

  return vIntersect, pIntersect, nil
}

func (m *HandshakeVersionMismatch) Error() error {
  // TODO: format the alternative possible versions
  return errors.New("Handshake version mismatch")
}

func (m *HandshakeDecodeError) Error() error {
  return errors.New("HandshakeDecodeError: " + m.Reason)
}

func (m *HandshakeRefused) Error() error {
  return errors.New("HandshakeRefused: " + m.Reason)
}
