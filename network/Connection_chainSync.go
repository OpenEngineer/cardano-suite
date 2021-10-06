package network

import (
  "errors"
  "fmt"

  "github.com/christianschmitz/cardano-suite/common"
)

type ChainSyncState int

const (
  CHAIN_SYNC_IDLE       ChainSyncState = iota
  CHAIN_SYNC_INTERSECT
  CHAIN_SYNC_CAN_AWAIT
  CHAIN_SYNC_MUST_REPLY
  CHAIN_SYNC_DONE
)

func (c *Connection) handleChainSyncSegment(b []byte) error {
  msgs_, err := common.AppendToCBORBuffer(c.chainSyncBuffer, b, MAX_PACKET_BUFFER_SIZE)
  if err != nil {
    return err
  }

  if msgs_ == nil {
    return nil
  }

  for _, msg_ := range msgs_ {
    msg, err := ChainSyncMessageFromUntyped(msg_)
    if err != nil {
      return err
    }
    
    c.chainSyncIngress <- msg
  }

  return nil
}

func (c *Connection) listenChainSync() {
  go func() {
    fmt.Println("[INFO] chainSync thread started")

    Outer:
    for {
      var msg_ ChainSyncMessage
      select {
      case msg_ = <- c.chainSyncIngress:
      case <- c.chainSyncCloseIngress:
        break Outer
      }

      switch msg := msg_.(type) {
      case *ChainSyncAwaitReply:
        switch c.chainSyncState {
        case CHAIN_SYNC_CAN_AWAIT:
          c.chainSyncState = CHAIN_SYNC_MUST_REPLY
        default:
          c.chainSyncErrorEgress <- errors.New("wrong state for receiving ChainSyncAwaitReply")
          <- c.chainSyncCloseIngress
          break Outer
        }
      case *ChainSyncRollForward:
        switch c.chainSyncState {
        case CHAIN_SYNC_CAN_AWAIT, CHAIN_SYNC_MUST_REPLY:
          if err := c.state.RollForward(c.nodeId, msg.Header.Header); err != nil {
            c.chainSyncErrorEgress <- err
            <- c.chainSyncCloseIngress
            break Outer
          }
          c.chainSyncState = CHAIN_SYNC_IDLE
        default:
          c.chainSyncErrorEgress <-errors.New("wrong state for receiving ChainSyncRollForwardMessage") 
          <- c.chainSyncCloseIngress
          break Outer
        }
      case *ChainSyncRollBackward:
        switch c.chainSyncState {
        case CHAIN_SYNC_CAN_AWAIT, CHAIN_SYNC_MUST_REPLY:
          if err := c.state.RollBackward(c.nodeId, msg.Point); err != nil {
            c.chainSyncErrorEgress <- err
            <- c.chainSyncCloseIngress
            break Outer
          }

          c.chainSyncState = CHAIN_SYNC_IDLE
        default:
          c.chainSyncErrorEgress <- errors.New("wrong state for receiving ChainSyncRollbackward")
          <- c.chainSyncCloseIngress
          break Outer
        }
      case *ChainSyncIntersectFound:
        if c.chainSyncState != CHAIN_SYNC_INTERSECT {
          c.chainSyncErrorEgress <- errors.New("wrong state for receiving ChainSyncIntersectFound")
          <- c.chainSyncCloseIngress
          break Outer
        }

        if err := c.state.SetNodeIntersect(c.nodeId, msg.Point); err != nil {
          c.chainSyncErrorEgress <- err
          <- c.chainSyncCloseIngress
          break Outer
        }

        c.chainSyncState = CHAIN_SYNC_IDLE
      case *ChainSyncIntersectNotFound:
        if c.chainSyncState != CHAIN_SYNC_INTERSECT {
          c.chainSyncErrorEgress <- errors.New("wrong state for receiving ChainSyncIntersectFound")
          <- c.chainSyncCloseIngress
          break Outer
        }

        // TODO: suggest other intersect
        panic("not yet implemented")

        c.chainSyncState = CHAIN_SYNC_IDLE
      default:
        c.chainSyncErrorEgress <- errors.New("unexpected message received during ChainSync")
        <- c.chainSyncCloseIngress
        break Outer
      }
    }

    c.chainSyncCloseEgress <- true
  }()
}
