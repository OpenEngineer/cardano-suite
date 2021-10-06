package cardano

import (
  "bytes"
  "encoding/binary"
  "errors"
  "fmt"
  "os"
  "strconv"
  "time"
)

// can't use automatic marshal/unmarshal tools due to weird IsResponder bit
type SegmentHeader struct {
  Timestamp    uint32
  IsResponder  bool
  ProtocolId   uint16
  PayloadSize  uint16
}

func NewSegmentHeader(isResponder bool, protocolId int, payloadSize int) *SegmentHeader {
  h := NewTimestampedSegmentHeader(0, isResponder, protocolId, payloadSize)

  h.SetTimeNow()

  return h
}

func NewTimestampedSegmentHeader(timestamp int64, isResponder bool, protocolId int, payloadSize int) *SegmentHeader {
  t32 := uint32(timestamp)

  return &SegmentHeader{
    t32,
    isResponder,
    uint16(protocolId),
    uint16(payloadSize),
  }
}

func (h *SegmentHeader) ToBytes() []byte {
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

func (h *SegmentHeader) SetTimeNow() {
  t64 := time.Now().UTC().Unix()*1000000

  h.Timestamp = uint32(t64)
}

func SegmentHeaderFromBytes(b []byte) (*SegmentHeader, []byte, error) {
  if len(b) < 8 {
    return nil, nil, errors.New("packet too small to be able to read SegmentHeader (" + strconv.Itoa(len(b)) + ")")
  }

  buf := bytes.NewBuffer(b)

  h := &SegmentHeader{}

  if err := binary.Read(buf, binary.BigEndian, &(h.Timestamp)); err != nil {
    return nil, nil, err
  }

  var pId uint16 
  if err := binary.Read(buf, binary.BigEndian, &pId); err != nil {
    return nil, nil, err
  }

  h.IsResponder = (pId & 0x8000) > 0
  h.ProtocolId = pId & (0x8000 - 1)

  if err := binary.Read(buf, binary.BigEndian, &h.PayloadSize); err != nil {
    return nil, nil, err
  }

  return h, b[8:], nil
}

func (h *SegmentHeader) Dump() {
  fmt.Fprintf(
    os.Stdout, 
    "SegmentHeader{\n  timestamp: %d,\n  isResponder: %t,\n  protocolId: %d,\n  payloadSize: %d\n}\n",
    h.Timestamp,
    h.IsResponder,
    h.ProtocolId,
    h.PayloadSize,
  )
}
