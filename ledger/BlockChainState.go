package ledger

import (
  "errors"
  "sync"

  "github.com/christianschmitz/cardano-suite/common"
)

// only the hot nodes can vote, so when the program starts the fork node lists are empty

type BlockChainState struct {
  PrevBlock     common.Hash  // not relevant when first block id of main chain is genesis

  LargestRollback int // count the largest rollback (i.e. fork switch), it is only sensible to download and process the blocks that are at least this deep on the main chain
  Forks []Fork        // the first fork is the main chain

  mutex *sync.RWMutex
}

type NodeState struct {
  node       common.Hash
  tipBlockId uint64
}

type Fork struct {
  ParentForkId    int    // -1 when on main chain, 0 to refer to main chain, other number when forking from other fork

  ParentBlockId   uint64 // irrelevant when prevBlockFork is -1

  Blocks []BlockHeader

  nodes  []NodeState   // hot nodes who voted for this fork, fork with largest number of votes becomes the main chain, not exported because we restart the client we want a fresh state of these
}

func NewGenesisBlockChainState() *BlockChainState {
  return &BlockChainState{
    common.NilHash(),
    0,
    []Fork{
      Fork{
        -1, 0, []BlockHeader{NewBlockHeader(NewGenesisHash(), 0, 0, common.HashFromString(""))}, []NodeState{},
      },
    },
    &sync.RWMutex{},
  }
}

func (f Fork) Len() int {
  return len(f.Blocks)
}

// returns -1 if not found
func (b *BlockChainState) findNodeFork(nodeId common.Hash) int {
  for i, f := range b.Forks {
    for _, n := range f.nodes {
      if n.node.Eq(nodeId) {
        return i
      }
    }
  }

  return -1
}

func (b *BlockChainState) lastBlockHeader(fId int) BlockHeader {
  if fId < 0 || fId >= len(b.Forks) {
    panic("fork id out of range")
  }

  lst := b.Forks[fId].Blocks

  return lst[len(lst)-1]
}

// if fId==-1 then all forks are searched
// searches from the end
func (b *BlockChainState) findBlock(fId int, headerHash common.Hash) (uint64, int) {
  for i, f := range b.Forks {
    if fId == -1 || fId == i {

      for bId_ := len(f.Blocks) - 1; bId_ >= 0; bId_-- {
        bl := f.Blocks[bId_]

        var bId uint64 = uint64(bId_)
        if f.ParentForkId != -1 {
          bId += f.ParentBlockId + 1
        }

        hh := bl.HeaderHash

        if hh.Eq(headerHash) {
          return bId, i
        }
      }
    }
  }

  return 0, -1
}

func (b *BlockChainState) findPrevBlock(fId int, h ChainHash) (uint64, int) {
  for i, f := range b.Forks {
    if fId == -1 || fId == i {

      for bId_ := len(f.Blocks) - 1; bId_ >= 0; bId_-- {
        bl := f.Blocks[bId_]

        var bId uint64 = uint64(bId_)
        if f.ParentForkId != -1 {
          bId += f.ParentBlockId + 1
        }

        pbh := NewChainHash(bl.HeaderHash)

        if pbh.Eq(h) {
          return bId, i
        }
      }

      if h.IsGenesis && f.ParentForkId == -1 {
        return 0, i
      }
    }
  }

  return 0, -1
}

func (b *BlockChainState) createFork(parentForkId int, parentBlockId uint64) int {
  fIdNew := len(b.Forks)

  b.Forks = append(b.Forks, Fork{
    parentForkId, 
    parentBlockId,
    []BlockHeader{},
    []NodeState{},
  })

  return fIdNew
}

func (b *BlockChainState) appendToFork(fId int, header *BlockHeader) error {
  last := b.lastBlockHeader(fId)

  if last.BlockId != header.BlockId - 1 {
    return errors.New("blockId doesnt follow last")
  }

  if last.SlotId >= header.SlotId {
    return errors.New("slotId not increasing")
  }

  if err := header.VerifyHash(); err != nil {
    return err
  }

  if !header.PrevBlock.Eq(NewChainHash(last.HeaderHash)) {
    return errors.New("chain hash not following")
  }

  b.Forks[fId].Blocks = append(b.Forks[fId].Blocks, *header)

  return nil
}

// returns an error if an inconsistency is found
// the connection with that node will probably be terminated by the caller in that case
// 
func (b *BlockChainState) RollForward(nodeId common.Hash, header *BlockHeader) error {
  b.mutex.Lock()

  defer b.mutex.Unlock()

  fId := b.findNodeFork(nodeId)

  if fId < 0 {
    // node not yet found
    _, fId = b.findPrevBlock(-1, header.PrevBlock)

    if fId < 0 {
      return errors.New("preceding block not found")
    }

  }

  f := b.Forks[fId]

  chainFound := false

  nBlocks := f.Len()

  for i := nBlocks - 1; i >= 0; i-- {
    var ch ChainHash

    if i == -1 {
      ch = NewGenesisHash()
    } else {
      ch = NewChainHash(f.Blocks[i].HeaderHash)
    }

    if ch.Eq(header.PrevBlock) {
      // ok, the new header chains this fork
      chainFound = true

      if i == nBlocks - 1 {
        if err := b.appendToFork(fId, header); err != nil {
          return err
        }
      } else if !f.Blocks[i+1].HeaderHash.Eq(header.HeaderHash) {
        fIdNew := b.createFork(fId, f.Blocks[i].BlockId)

        if err := b.appendToFork(fIdNew, header); err != nil {
          return err
        }

        fId = fIdNew
      }

      break
    } 
  }

  if !chainFound {
    return errors.New("preceding block not found")
  }

  b.updateNodeState(fId, nodeId, header.BlockId)

  return nil
}

func (b *BlockChainState) RollBackward(nodeId common.Hash, point *Point) error {
  b.mutex.Lock()

  defer b.mutex.Unlock()

  fId := b.findNodeFork(nodeId)

  var bId uint64

  if fId < 0 {
    bId, fId = b.findBlock(-1, point.HeaderHash)

    if fId < 0 {
      return errors.New("block not found")
    }
  } else {
    bId, _ = b.findBlock(fId, point.HeaderHash)
  }

  fIdNew := b.createFork(fId, bId)

  b.updateNodeState(fIdNew, nodeId, bId)

  return nil
}

func (b *BlockChainState) updateNodeState(fId int, nodeId common.Hash, bId uint64) {
  for i, f := range b.Forks {
    if i != fId {
      for j, n := range f.nodes {
        if n.node.Eq(nodeId) {
          b.Forks[i].nodes = append(f.nodes[0:j], f.nodes[j+1:]...)
          break
        }
      }
    } else {
      found := false
      for j, n := range f.nodes {
        if n.node.Eq(nodeId) {
          found = true
          b.Forks[i].nodes[j] = NodeState{nodeId, bId}
          break
        }
      }

      if !found {
        b.Forks[i].nodes = append(b.Forks[i].nodes, NodeState{nodeId, bId})
      }
    }
  }
}

func (b *BlockChainState) SetNodeIntersect(nodeId common.Hash, point *Point) error {
  bId, fId := b.findBlock(-1, point.HeaderHash)

  b.updateNodeState(fId, nodeId, bId)

  return nil
}
