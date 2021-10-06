package network

import (
  "bytes"
  "fmt"
  "strconv"
  "time"
)

const RECOMMENDED_TCP_PAYLOAD_SIZE = 2900

func (c *Connection) startMux() {
  go func() {
    Outer:
    for {
      // a separate select for each mini-protocol

      select {
      case <- c.muxCloseIngress:
        break Outer
      case msg := <- c.handshakeEgress:
        c.pushOutgoingSegments(SegmentsFromBytes(false, HANDSHAKE, HandshakeMessageToCBOR(msg)))
      default:
      }

      if c.versionId == 0 { // handshake not yet finished, dont let the other mini-protocols send messages
        // dont block, just wait a bit
        time.Sleep(10*time.Millisecond)
      }  else {
        time.Sleep(time.Millisecond) // to avoid an endless cycle of unavailable messages

        select {
        case msg := <- c.chainSyncEgress:
          c.pushOutgoingSegments(SegmentsFromBytes(false, CHAIN_SYNC, ChainSyncMessageToCBOR(msg)))
        default:
        }

        select {
        case msg := <- c.blockFetchEgress:
          c.pushOutgoingSegments(SegmentsFromBytes(false, BLOCK_FETCH, BlockFetchMessageToCBOR(msg)))
        default:
        }

        select {
        case _ = <- c.txSubmissionEgress:
          panic("not yet implemented")
          //c.pushOutgoingSegments(SegmentsFromBytes(false, TX_SUBMISSION, TxSubmissionToCBOR(msg)))
        default:
        }

        select {
        case _ = <- c.keepAliveEgress:
          panic("not yet implemented")
          //c.pushOutgoingSegments(SegmentsFromBytes(false, KEEP_ALIVE, KeepAliveToCBOR(msg)))
        default:
        }
      }


      // intersperse segments from each mini-protocol
      ss := c.popOutgoingSegments()

      for len(ss) > 0 {
        buf := &bytes.Buffer{}

        // fill the tcp payload buffer until we are one segment above the recommended size
        for buf.Len() < RECOMMENDED_TCP_PAYLOAD_SIZE && len(ss) > 0 {
          ss[0].header.SetTimeNow()

          buf.Write(ss[0].ToBytes())

          ss = ss[1:]
        }

        data := buf.Bytes()

        _, err := c.conn.Write(data)
        if err != nil {
          c.muxErrorEgress <- err
          <- c.muxCloseIngress
          break Outer
        }
      }
    }

    c.muxCloseEgress <- true
  }()
}

func (c *Connection) pushOutgoingSegments(ss []*Segment) {
  if ss == nil || len(ss) == 0 {
    return
  }

  pId := ss[0].header.ProtocolId

  for _, ss_ := range ss {
    fmt.Println("[INFO] sending segment", *ss_, "(protocol " + strconv.Itoa(int((*ss_).header.ProtocolId)) + ")")
  }

  if c.outgoingSegments[pId] == nil {
    c.outgoingSegments[pId] = ss
  } else {
    c.outgoingSegments[pId] = append(c.outgoingSegments[pId], ss...)
  }
}

func (c *Connection) popOutgoingSegments() []*Segment {
  res := make([]*Segment, 0)

  maybeHavePendingSegments := true

  for maybeHavePendingSegments {
    maybeHavePendingSegments = false

    for pId, ss := range c.outgoingSegments {
      if ss != nil && len(ss) > 0 {
        res = append(res, ss[0])

        c.outgoingSegments[pId] = ss[1:]

        if len(ss) > 1 {
          maybeHavePendingSegments = true
        }
      }
    }
  }

  return res
}
