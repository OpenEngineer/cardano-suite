package network

import (
  "errors"
  "fmt"
  "reflect"
  "strconv"

  "github.com/christianschmitz/cardano-suite/common"
)

type HandshakeState int

const (
  HANDSHAKE_PROPOSE HandshakeState = iota
  HANDSHAKE_CONFIRM
  HANDSHAKE_DONE
)

func (c *Connection) HandshakeState() HandshakeState {
  return c.handshakeState
}

func (c *Connection) handleHandshakeSegment(b []byte) error {
  msgs_, err := common.AppendToCBORBuffer(c.handshakeBuffer, b, MAX_PACKET_BUFFER_SIZE)
  if err != nil {
    return err
  }

  if msgs_ == nil {
    return nil
  }

  for _, msg_ := range msgs_ {
    msg, err := HandshakeMessageFromUntyped(msg_)
    if err != nil {
      return err
    }
    
    c.handshakeIngress <- msg
  }

  return nil
}

// start the handshake thread
func (c *Connection) listenHandshake() {
  go func() {
    fmt.Println("[INFO] handshake thread started")

    Outer:
    for {
      var msg_ HandshakeMessage

      // when trying to close this thread using the handshakeCloseIngress in a non-blocking way. If it is blocked, then a msg is being processed, which will always end in handshakeCloseEgress of handshakeErrorEgress channel being used
      select {
      case msg_ = <- c.handshakeIngress:
      case <- c.handshakeCloseIngress:
        break Outer
      }

      // type assertions are used here so it is clear how the state machine evolves
      switch msg := msg_.(type) {
      case *HandshakeProposeVersions:
        if c.debug {
          fmt.Println("[INFO] received HandshakeProposeVersions")
        }

        v, p, err := msg.Intersect(DefaultProtocolVersions(c.magic))
        if err != nil {
          c.handshakeErrorEgress <- err // dont actually quit upon error, quit is handled by the close
        } else {
          c.finalizeHandshake(v, p)

          c.handshakeState = HANDSHAKE_DONE
        }

        <- c.handshakeCloseIngress
        break Outer
      case *HandshakeAcceptVersion:
        if c.debug {
          fmt.Println("[INFO] received HandshakeAcceptVersion (" + strconv.Itoa(msg.Version) +  ")")
        }

        c.finalizeHandshake(msg.Version, msg.Params)

        c.handshakeState = HANDSHAKE_DONE

        <- c.handshakeCloseIngress
        break Outer
      case HandshakeError:
        c.handshakeState = HANDSHAKE_DONE

        c.handshakeErrorEgress <- msg.Error()
        <- c.handshakeCloseIngress
        break Outer
      default:
        c.handshakeErrorEgress <- errors.New("unhandled HandshakeMessage type " + reflect.TypeOf(msg_).String())
        <- c.handshakeCloseIngress // dont just exit, wait for calling thread to issue close
        break Outer
      }
    }

    fmt.Println("[INFO] handshake thread closed")
    c.handshakeCloseEgress <- true // closed prematurely
  }()
}

func (c *Connection) initHandshake() {
  c.handshakeEgress <- NewHandshakeProposeVersions(c.magic)

  c.handshakeState = HANDSHAKE_CONFIRM
}

func (c *Connection) finalizeHandshake(v int, p interface{}) {
  c.versionId, c.versionParams = v, p

  c.handshakeVersionEgress <- v
}
