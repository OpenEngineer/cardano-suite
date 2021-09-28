package cardano

import (
  "bytes"
  "encoding/binary"
  "errors"
  "fmt"
  "os"
  "time"
)

// can't use automatic marshal/unmarshal tools due to weird IsResponder bit
type MuxHeader struct {
  Timestamp    uint32
  IsResponder  bool
  ProtocolId   uint16
  PayloadSize  uint16
}

func NewMuxHeader(isResponder bool, protocolId int, payloadSize int) *MuxHeader {
  // according to ouroboros-network/network-mux/src/Network/Mux/Time.hs the epoch can be arbitrary and the timestamp is only used for differences wrt. the previous packets
  t64 := time.Now().UTC().Unix()*1000000

  return NewTimestampedMuxHeader(t64, isResponder, protocolId, payloadSize)
}

func NewTimestampedMuxHeader(timestamp int64, isResponder bool, protocolId int, payloadSize int) *MuxHeader {
  t32 := uint32(timestamp)

  return &MuxHeader{
    t32,
    isResponder,
    uint16(protocolId),
    uint16(payloadSize),
  }
}

func (h *MuxHeader) ToBytes() []byte {
  buf := bytes.NewBuffer(nil)

  if err := binary.Write(buf, binary.BigEndian, h.Timestamp); err != nil {
    panic(err)
  }

  pId := h.ProtocolId

  if h.IsResponder {
    pId = pId | 0x8000
  }

  if err := binary.Write(buf, binary.BigEndian, pId); err != nil {
    panic(err)
  }

  if err := binary.Write(buf, binary.BigEndian, h.PayloadSize); err != nil {
    panic(err)
  }

  return buf.Bytes()
}

func MuxHeaderFromBytes(b []byte) (*MuxHeader, error) {
  if len(b) < 8 {
    return nil, errors.New("packet too small to be able to read MuxHeader")
  }

  buf := bytes.NewBuffer(b)

  h := MuxHeader{}

  if err := binary.Read(buf, binary.BigEndian, &(h.Timestamp)); err != nil {
    return nil, err
  }

  var pId uint16 
  if err := binary.Read(buf, binary.BigEndian, &pId); err != nil {
    return nil, err
  }

  h.IsResponder = (pId & 0x8000) > 0
  h.ProtocolId = pId & (0x8000 - 1)

  if err := binary.Read(buf, binary.BigEndian, &h.PayloadSize); err != nil {
    return nil, err
  }

  return &h, nil
}

func (h *MuxHeader) Dump() {
  fmt.Fprintf(
    os.Stdout, 
    "MuxHeader{\n  timestamp: %d,\n  isResponder: %t,\n  protocolId: %d,\n  payloadSize: %d\n}\n",
    h.Timestamp,
    h.IsResponder,
    h.ProtocolId,
    h.PayloadSize,
  )
}
