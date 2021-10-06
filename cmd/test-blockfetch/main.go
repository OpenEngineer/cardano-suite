package main

import (
  "errors"
  "fmt"
  "os"
  "strconv"

  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
  "github.com/christianschmitz/cardano-suite/network"
)

// TODO: make this a proper utility, whereby the several servers are tried, and the blockchain is tried
// some kind of channel is needed to wait for the result
// also proper error management

const (
  TEST_SERVER  = "relays-new.cardano-testnet.iohkdev.io" 
  DEFAULT_PORT = "3001"
  MAGIC        = network.TESTNET_MAGIC
)

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func mainInternal() error {
  args := os.Args[1:]
  
  if len(args) != 2 {
    return errors.New("Usage: test-blockfetch slotId headerHash")
  }

  slotId, err := strconv.ParseUint(args[0], 10, 64)
  if err != nil {
    return err
  }

  headerHash, err := common.HashFromHex(args[1])
  if err != nil {
    return err
  }
  

  state := ledger.NewGenesisBlockChainState()

  conn, err := network.NewConnection(MAGIC, TEST_SERVER, DEFAULT_PORT, state, true)
  if err != nil {
    return err
  }

  bl, err := conn.FetchSingleBlock(slotId, headerHash)
  if err != nil {
    return err
  }

  fmt.Println(common.ReflectUntyped(bl.Untyped))

  return nil
}
