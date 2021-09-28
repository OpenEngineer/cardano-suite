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

func parseArgsAsBytes(args []string) ([]byte, error) {
  b := make([]byte, len(args))

  for i, arg := range args {
    v, err := strconv.ParseInt(arg, 16, 64)
    if err != nil {
      return nil, errors.New("failed to parse byte " + strconv.Itoa(i))
    }

    b[i] = byte(v)
  }

  return b, nil
}

func mainInternal() error {
  args := os.Args[1:]

  if len(args) == 8 {
    return testBytes(args)
  } else if len(args) == 0 {
    return testNow()
  } else if len(args) == 1 {
    return testTimestamped(args[0])
  } else {
    return errors.New("unhandled number of args")
  }
}

func testBytes(args []string) error {
  b, err := parseArgsAsBytes(args)
  if err != nil {
    return err
  }

  fmt.Println("Integers:", b)

  h, err := cardano.MuxHeaderFromBytes(b)
  if err != nil {
    return err
  }

  h.Dump()

  bOrig := h.ToBytes()
  fmt.Println("Integers':", bOrig)

  return nil
}

func testNow() error {
  // using now
  h := cardano.NewMuxHeader(false, 0, 45)

  return testHeader(h)
}

func testTimestamped(arg string) error {
  t64, err := strconv.ParseInt(arg, 10, 64)
  if err != nil {
    return err
  }

  h := cardano.NewTimestampedMuxHeader(t64, false, 0, 45)

  return testHeader(h)
}

func testHeader(h *cardano.MuxHeader) error {
  h.Dump()

  b := h.ToBytes()
  fmt.Println("Integers':", b)

  hOrig, err := cardano.MuxHeaderFromBytes(b)
  if err != nil {
    return err
  }

  hOrig.Dump()

  return nil
}
