package network

import (
  "errors"
  "fmt"

  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

type BlockFetchState int

const (
  BLOCK_FETCH_IDLE BlockFetchState = iota
  BLOCK_FETCH_BUSY 
  BLOCK_FETCH_STREAMING
  BLOCK_FETCH_DONE
)

func (c *Connection) handleBlockFetchSegment(b []byte) error {
  msgs_, err := common.AppendToCBORBuffer(c.blockFetchBuffer, b, MAX_PACKET_BUFFER_SIZE)
  if err != nil {
    return err
  }

  if msgs_ == nil {
    return nil
  }

  for _, msg_ := range msgs_ {
    msg, err := BlockFetchMessageFromUntyped(msg_)
    if err != nil {
      return err
    }
    
    c.blockFetchIngress <- msg
  }

  return nil
}

func (c *Connection) listenBlockFetch() {
  go func() {
    fmt.Println("[INFO] blockFetch thread started")
    Outer:
    for {
      var msg_ BlockFetchMessage
      select {
      case msg_ = <- c.blockFetchIngress:
      case <- c.blockFetchCloseIngress:
        c.blockFetchState = BLOCK_FETCH_DONE
        break Outer
      }

      switch msg := msg_.(type) {
      case *BlockFetchClientDone:

        c.blockFetchState = BLOCK_FETCH_DONE

        <- c.blockFetchIngress
        break Outer
      case *BlockFetchStartBatch:
        if c.blockFetchState != BLOCK_FETCH_BUSY {
          c.blockFetchErrorEgress <- errors.New("wrong state for receiving BlockFetchNoBlocks")
          <- c.blockFetchIngress
          break Outer
        }

        fmt.Println("[INFO] started block fetch streaming")

        c.blockFetchState = BLOCK_FETCH_STREAMING
      case *BlockFetchNoBlocks:
        if c.blockFetchState != BLOCK_FETCH_BUSY {
          c.blockFetchErrorEgress <- errors.New("wrong state for receiving BlockFetchNoBlocks")
          <- c.blockFetchIngress
          break Outer
        }


        c.blockFetchState = BLOCK_FETCH_IDLE

        if c.BlockFetchEventEgress != nil {
          c.blockFetchErrorEgress <- errors.New("no blocks found in specified range")
          <- c.blockFetchIngress
          break Outer
        } else {
          fmt.Println("[INFO] no blocks found in specified range")
        }
      case *BlockFetchBlock:
        if c.blockFetchState != BLOCK_FETCH_STREAMING {
          c.blockFetchErrorEgress <- errors.New("wrong state for receiving block")
          <- c.blockFetchIngress
          break Outer
        }

        // TODO: call a callback depending on the slotNo

        if c.BlockFetchEventEgress != nil {
          c.BlockFetchEventEgress <- msg.Block.Block
        } else {
          str:= common.ReflectUntyped(msg.Block.Block.Untyped)
          fmt.Println("#received block:", str)
        }
      case *BlockFetchBatchDone:
        if c.blockFetchState != BLOCK_FETCH_STREAMING {
          c.blockFetchErrorEgress <- errors.New("wrong state for receiving BlockFetchBatchDone")
          <- c.blockFetchIngress
          break Outer
        }

        fmt.Println("[INFO] block batch done")

        c.blockFetchState = BLOCK_FETCH_IDLE

        if c.BlockFetchEventEgress != nil {
          <- c.blockFetchIngress
          break Outer
        }
      default:
        c.blockFetchErrorEgress <- errors.New("unexpected BlockFetchMessage type received from server")
        <- c.blockFetchIngress
        break Outer
      }
    }

    c.blockFetchCloseEgress <- true
  }()
}

// TODO: register a channel as a callback
func (c *Connection) RequestBlockFetchRange(msg *BlockFetchRequestRange) {
  // dont block
  go func() {
    c.blockFetchState = BLOCK_FETCH_BUSY

    c.blockFetchEgress <- msg
  }()
}

// can be called synchronously from the main thread
// starts the listen thread
func (c *Connection) FetchSingleBlock(slotId uint64, headerHash common.Hash) (*ledger.Block, error) {
  errCh := make(chan error)

  c.BlockFetchEventEgress = make(chan *ledger.Block)

  c.ListenNohup(errCh)

  point := &ledger.Point{false, slotId, headerHash}

  msg := &BlockFetchRequestRange{point, point}

  c.RequestBlockFetchRange(msg)

  select {
  case err := <- errCh:
    return nil, err
  case bl := <- c.BlockFetchEventEgress:
    // dont bother closing the threads properly
    return bl, nil
  }
}
