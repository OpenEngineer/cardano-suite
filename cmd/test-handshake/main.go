package main

import (
  "errors"
  "fmt"
  "io"
  "net"
  "os"
  "strings"
  "time"

  cardano "github.com/christianschmitz/cardano-suite"
)

// About
//  This utility intiates a tcp connection with a public cardano-node, 
//  and subsequently runs the Handshake protocol, after which the connection is closed.
//
//  Running this test multiple times in a row will fail sometimes due to DoS prevention.
//  Leave at least 20s between tries.

// testnet server example:
//  TODO: resolve dns into multiple ip addresses
const TEST_SERVER = "relays-new.cardano-testnet.iohkdev.io:3001" 
const MAGIC = cardano.TESTNET_MAGIC

// mainnet server example:
//const TEST_SERVER = "50.18.237.14:3001"
//const MAGIC = cardano.MAINNET_MAGIC


type Context struct {
  version int
}

func (ctx *Context) Version() int {
  return ctx.version
}

func (ctx *Context) SetVersion(v int) {
  ctx.version = v
}

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func send(conn io.Writer, msg cardano.HandshakeMessage) error {
  bPayload := cardano.HandshakeMessageToCBOR(msg)

  header := cardano.NewMuxHeader(false, 0, len(bPayload))

  bHeader := header.ToBytes()

  b := append(bHeader, bPayload...)

  _, err := conn.Write(b)

  return err
}

func receive(conn io.Reader) (cardano.HandshakeMessage, error) {
  bHeader := make([]byte, 8)

  nHeader, err := conn.Read(bHeader)
  if err != nil {
    if nHeader != 0 {
      fmt.Println("[INFO] EOF error while reading", nHeader, " bytes")
    } else {
      fmt.Println("[INFO] Nothing received\n")
    }

    return nil, err
  }

  header, err := cardano.MuxHeaderFromBytes(bHeader)
  if err != nil {
    return nil, err
  }

  header.Dump()

  bPayload := make([]byte, header.PayloadSize)

  nPayload, err := conn.Read(bPayload)
  if err != nil {
    return nil, err
  }

  if nPayload != len(bPayload) {
    return nil, errors.New("worng payload size")
  }

  fmt.Println("[INFO] Received", nPayload, "bytes:", bPayload)

  return cardano.HandshakeMessageFromCBOR(bPayload)
}

func pickUrl(url string) (string, error) {
  hostPort := strings.Split(url, ":")
  if len(hostPort) != 2 {
    return "", errors.New("expected host:port in url")
  }

  host := hostPort[0]
  port := hostPort[1]
  
  addrs, err := net.LookupHost(host)
  if err != nil {
    return "", err
  }

  if len(addrs) == 0 {
    return "", errors.New("no dns records found")
  }

  addr := addrs[len(addrs)-1]

  if len(addrs) != 1 {
    fmt.Println("[INFO] Picked host", addr)
  }

  return addr + ":" + port, nil
}

func mainInternal() error {
  args := os.Args[1:]

  if len(args) == 0 {
    args = []string{TEST_SERVER}
  } else if len(args) != 1 {
    return errors.New("Usage: test-handshake <host:port>\n")
  }

  url, err := pickUrl(args[0])
  if err != nil {
    return err
  }

  // first establish a tcp connection

  conn, err := net.Dial("tcp", url)
  if err != nil {
    return err
  }

  if err := conn.SetReadDeadline(time.Time{}); err != nil {
    return err
  }

  tcpConn, ok := conn.(*net.TCPConn)
  if !ok {
    return errors.New("not a tcp connection")
  }

  if err := tcpConn.SetLinger(0); err != nil {
    return err
  }

  fmt.Println("[INFO] Dial ok")

  ctx := &Context{-1}

  var msg cardano.HandshakeMessage = cardano.NewHandshakeProposeVersions(MAGIC)

  for i := 0; msg != nil; i++ {
    fmt.Println("[INFO] Sending msg", i, "...")

    send(conn, msg)

    fmt.Println("[INFO] Sent")

    fmt.Println("[INFO] Receiving resp to msg", i, "...")

    msgR, err := receive(conn)
    if err != nil {
      return err
    }

    msg, err = msgR.Respond(ctx)
    if err != nil {
      return err
    }
  }

  fmt.Println("Handshake version:", ctx.Version())

  return conn.Close()
}
