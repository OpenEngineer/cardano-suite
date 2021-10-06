package main

import (
  "errors"
  "fmt"
  "os"
  "strconv"

  cardano "github.com/christianschmitz/cardano-suite"
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

  prevHash, err := cardano.HashFromHex(args[0])
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

  bodyHash := cardano.HashFromString(args[3])

  blockHeader := cardano.BlockHeader{
    cardano.NilHash(),
    cardano.ChainHash{false, prevHash},
    slotId,
    blockId,
    bodyHash,
  }

  headerHash := blockHeader.ToHash()

  fmt.Println(headerHash.ToHex())

  return nil
}
