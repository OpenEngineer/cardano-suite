package network

import (
  "errors"
  "fmt"
  "regexp"
  "strconv"
)

const MAX_PACKET_BUFFER_SIZE = 5000

func isErrClosed(err error) bool {
  if err == nil {
    return false
  } else {
    re := regexp.MustCompile(`use of closed network connection`)

    return re.MatchString(err.Error())
  }
}

// turn bytes into complete segments
func (c *Connection) startDemux() {
  buf := make([]byte, MAX_PACKET_BUFFER_SIZE)

  go func() {
    Outer:
    for {
      n, err := c.conn.Read(buf)
      if isErrClosed(err) {
        break Outer
      } else if err != nil {
        c.demuxErrorEgress <- err // dont actually quit, this is handled by c.conn.Close()
        break Outer
      } else {
        b := buf[0:n]

        fmt.Println("[INFO] received", b)

        for ;len(b) > 0; {
          if c.haveCompleteSegment() {
            if err := c.handleIncomingSegment(); err != nil {
              c.demuxErrorEgress <- err
              break Outer
            }
          }

          if c.incomingSegmentHeader == nil {
            c.incomingSegmentHeader, b, err = SegmentHeaderFromBytes(b)
            if err != nil {
              c.demuxErrorEgress <- err
              break Outer
            }
          } 

          if len(b) > 0 {
            // appending to incoming segment
            b, err = c.appendToIncomingSegment(b)
            if err != nil {
              c.demuxErrorEgress <- err
              break Outer
            }
          }
        }

        if c.haveCompleteSegment() {
          if err := c.handleIncomingSegment(); err != nil {
            c.demuxErrorEgress <- err
            break Outer
          }
        }
      }
    }

    c.demuxCloseEgress <- true
  }()
}

func (c *Connection) haveCompleteSegment() bool {
  if c.incomingSegmentHeader != nil {
    return c.incomingSegmentSize() <= c.incomingSegmentBuffer.Len()
  } else {
    return false
  }
}

func (c *Connection) incomingSegmentSize() int {
  if c.incomingSegmentHeader == nil {
    return 0
  } else {
    return int(c.incomingSegmentHeader.PayloadSize)
  }
}

// returns remaining bytes (i.e. that belong to other segments)
func (c *Connection) appendToIncomingSegment(b []byte) ([]byte, error) {
  remSize := c.incomingSegmentSize() - c.incomingSegmentBuffer.Len()

  if remSize < 0 {
    return nil, errors.New("buffer error")
  }

  var (
    this []byte
    rest []byte
  )

  if remSize >= len(b) {
    this = b
    rest = make([]byte, 0)
  } else {
    this = b[0:remSize]
    rest = b[remSize:]
  }

  _, err := c.incomingSegmentBuffer.Write(this)
  if err != nil {
    return nil, err
  }

  return rest, nil
}

func (c *Connection) handleIncomingSegment() error {
  h := c.incomingSegmentHeader
  if h == nil {
    return nil
  }

  b := c.incomingSegmentBuffer.Bytes()

  c.incomingSegmentHeader = nil
  c.incomingSegmentBuffer.Reset()
  // the different mini-protocols might handle the bytes differently, so send the raw bytes instead of unmarshaled CBOR

  switch h.ProtocolId {
  case HANDSHAKE:
    return c.handleHandshakeSegment(b)
  case CHAIN_SYNC:
    if c.versionId == 0 {
      panic("handshake not yet complete")
    }
    return c.handleChainSyncSegment(b)
  case BLOCK_FETCH:
    if c.versionId == 0 {
      panic("handshake not yet complete")
    }
    return c.handleBlockFetchSegment(b)
  case TX_SUBMISSION:
    if c.versionId == 0 {
      panic("handshake not yet complete")
    }
    return c.handleTxSubmissionSegment(b)
  case KEEP_ALIVE:
    if c.versionId == 0 {
      panic("handshake not yet complete")
    }
    return c.handleKeepAliveSegment(b)
  default:
    return errors.New("unknown protocol id " + strconv.Itoa(int(h.ProtocolId)))
  }
}

func (c *Connection) handleTxSubmissionSegment(b []byte) error {
  return errors.New("not yet implemented")


  /*msgs_, err := appendToCBORBuffer(c.txSubmissionBuffer, b, MAX_PACKET_BUFFER_SIZE)
  if err != nil {
    return err
  }

  if msgs_ == nil {
    return nil
  }

  for _, msg_ := range msgs_ {
    msg, err := TxSubmissionMessageFromUntyped(msg_)
    if err != nil {
      return err
    }
    
    c.txSubmissionIngress <- msg
  }

  return nil*/
}

func (c *Connection) handleKeepAliveSegment(b []byte) error {
  return errors.New("not yet implemented")

  // each KeepAlive segment is a valid message
  // they are not split because they are very small

  /*msg, err := KeepAliveMessageFromCBOR(b)
  if err != nil {
    return err
  }

  c.keepAliveIngress <- msg

  return nil*/
}
