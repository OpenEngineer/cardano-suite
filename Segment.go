package cardano

const SEGMENT_HEADER_SIZE = 8
const MAX_SEGMENT_PAYLOAD_SIZE = 688
const MAX_SEGMENT_SIZE = MAX_SEGMENT_PAYLOAD_SIZE + SEGMENT_HEADER_SIZE

type Segment struct {
  header  *SegmentHeader
  payload []byte
}

func (s *Segment) ToBytes() []byte {
  hb := s.header.ToBytes()

  return append(hb, s.payload...)
}

func SegmentFromBytes(isResponder bool, protocolId int, b []byte) *Segment {
  return &Segment{NewSegmentHeader(isResponder, protocolId, len(b)), b}
}

// XXX: this is not as smart as the haskell implementation, but does it really matter? (I dont think so because the segment headers are tiny
func SegmentsFromBytes(isResponder bool, protocolId int, b []byte) []*Segment {
  res := make([]*Segment, 0)

  for i := 0; i < len(b); i += MAX_SEGMENT_PAYLOAD_SIZE {
    res = append(res, SegmentFromBytes(isResponder, protocolId, b[i:min(i+MAX_SEGMENT_PAYLOAD_SIZE, len(b))]))
  }

  return res
}
