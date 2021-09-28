package cardano

import (
  "errors"
)

const MAINNET_MAGIC = 764824073
const TESTNET_MAGIC = 1097911063

type HandshakeContext interface {
  SetVersion(v int)
  Version() int
}

// TODO: compare cbor results to a wireshark capture
type HandshakeMessage interface { // all messages of a given protocol need to have the same interface type
  Respond(ctx HandshakeContext) (HandshakeMessage, error) // when received from server, if the follow-up message is nil nothing is sent
}

type HandshakeProposeVersions struct {
  Versions map[int]interface{} // XXX: does cbor result need to be sorted?
}

type HandshakeAcceptVersion struct {
  Version int
  Param   interface{}
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

// NOTE: the primitive type-names below are different, because the builtin cast functions are not very robust
//    see ./cbor_primitive_types.go for the possible cast functions

//go:generate ./cbor-type HandshakeMessage *HandshakeProposeVersions *HandshakeAcceptVersion HandshakeRefuse
//go:generate ./cbor-type *HandshakeProposeVersions "Versions IntToInterfMap"
//go:generate ./cbor-type *HandshakeAcceptVersion   "Version Int, Param Interf"
//go:generate ./cbor-type HandshakeRefuse *HandshakeVersionMismatch *HandshakeDecodeError *HandshakeRefused
//go:generate ./cbor-type *HandshakeVersionMismatch "ValidVersions IntList"
//go:generate ./cbor-type *HandshakeDecodeError     "Version Int, Error String"
//go:generate ./cbor-type *HandshakeRefused         "Version Int, Error String"

func NewHandshakeProposeVersions(magic int) *HandshakeProposeVersions {
  m := &HandshakeProposeVersions{make(map[int]interface{})}

  m.Versions[1] = magic
  m.Versions[2] = magic
  m.Versions[3] = magic

  l4 := make([]interface{}, 2)
  l4[0] = magic
  l4[1] = false
  m.Versions[4] = l4

  l5 := make([]interface{}, 2)
  l5[0] = magic
  l5[1] = false
  m.Versions[5] = l5

  l6 := make([]interface{}, 2)
  l6[0] = magic
  l6[1] = false
  m.Versions[6] = l6

  return m
}

func (m *HandshakeProposeVersions) Respond(ctx HandshakeContext) (HandshakeMessage, error) {
  return nil, errors.New("unhandled ProposeVersions response from server")
}

func (m *HandshakeAcceptVersion) Respond(ctx HandshakeContext) (HandshakeMessage, error) {
  ctx.SetVersion(m.Version)

  return nil, nil
}

func (m *HandshakeVersionMismatch) Respond(ctx HandshakeContext) (HandshakeMessage, error) {
  return nil, errors.New("Handshake version mismatch")
}

func (m *HandshakeDecodeError) Respond(ctx HandshakeContext) (HandshakeMessage, error) {
  return nil, errors.New("HandshakeDecodeError: " + m.Error)
}

func (m *HandshakeRefused) Respond(ctx HandshakeContext) (HandshakeMessage, error) {
  return nil, errors.New("HandshakeRefused: " + m.Error)
}
