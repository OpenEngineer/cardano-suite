package main

import (
  "errors"
  "fmt"
  "os"
  "strconv"

  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func mainInternal() error {
  args := os.Args[1:]

  if len(args) != 4 {
    return errors.New("expected 4 args (prevBlockHash, slotId, blockId, body)")
  }

  prevHash, err := common.HashFromHex(args[0])
  if err != nil {
    return err
  }

  slotId, err := strconv.ParseUint(args[1], 10, 64)
  if err != nil {
    return err
  }

  blockId, err := strconv.ParseUint(args[2], 10, 64)
  if err != nil {
    return err
  }

  bodyHash := common.HashFromString(args[3])

  blockHeader := ledger.BlockHeader{
    common.NilHash(),
    ledger.ChainHash{false, prevHash},
    slotId,
    blockId,
    bodyHash,
  }

  headerHash := blockHeader.ToHash()

  fmt.Println(headerHash.ToHex())

  return nil
}
