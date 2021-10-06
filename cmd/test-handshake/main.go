package main

import (
  "errors"
  "fmt"
  "os"
  "time"

  "github.com/christianschmitz/cardano-suite/ledger"
  "github.com/christianschmitz/cardano-suite/network"
)

// About
//  This utility intiates a tcp connection with a public cardano-node, 
//  and subsequently runs the Handshake protocol, after which the connection is closed.
//
//  Running this test multiple times in a row will fail sometimes due to DoS prevention.
//  Leave at least 20s between tries.

// testnet server example:
//  TODO: resolve dns into multiple ip addresses
const (
  TEST_SERVER  = "relays-new.cardano-testnet.iohkdev.io" 
  DEFAULT_PORT = "3001"
  MAGIC        = network.TESTNET_MAGIC
  // MAGIC = network.MAINNET_MAGIC
)

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func mainInternal() error {
  args := os.Args[1:]

  if len(args) == 0 {
    args = []string{TEST_SERVER, DEFAULT_PORT}
  } else if len(args) != 2 && len(args) != 1 {
    return errors.New("Usage: test-handshake host [port]\n")
  }

  if len(args) == 1 {
    args = append(args, DEFAULT_PORT)
  }

  host := args[0]
  port := args[1]

  state := ledger.NewGenesisBlockChainState()

  conn, err := network.NewConnection(MAGIC, host, port, state, true)
  if err != nil {
    return err
  }

  errCh := make(chan error)

  conn.ListenNohup(errCh)

  for {
    select {
    case err := <- errCh:
      return err
    default:
      if conn.Version() == 0 {
        time.Sleep(time.Millisecond*100)
      } else {
        fmt.Println("Negotiated version: ", conn.Version())
        conn.Close()
        return nil
      }
    }
  }

  return nil
}
