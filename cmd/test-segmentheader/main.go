package main

import (
  "errors"
  "fmt"
  "os"
  "regexp"
  "strconv"

  cardano "github.com/christianschmitz/cardano-suite"
)

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func parseArgsAsBytes(args []string, base int) ([]byte, error) {
  b := make([]byte, len(args))

  for i, arg := range args {
    v, err := strconv.ParseInt(arg, base, 64)
    if err != nil {
      return nil, errors.New("failed to parse byte " + strconv.Itoa(i))
    }

    b[i] = byte(v)
  }

  return b, nil
}

func allInts(args []string) bool {
  re := regexp.MustCompile(`^[0-9]*$`)

  for _, arg := range args {
    if !re.MatchString(arg) {
      return false
    }
  }

  return true
}

func mainInternal() error {
  args := os.Args[1:]

  if len(args) == 8 {
    var (
      b []byte
      err error
    )

    if allInts(args) {
      b, err = parseArgsAsBytes(args, 10)
      if err != nil {
        return err
      }
    } else {
      b, err = parseArgsAsBytes(args, 16)
      if err != nil {
        return err
      }
    }

    return testBytes(b)
  } else if len(args) == 0 {
    return testNow()
  } else if len(args) == 1 {
    return testTimestamped(args[0])
  } else {
    return errors.New("unhandled number of args")
  }
}

func testBytes(b []byte) error {
  fmt.Println("Integers:", b)

  h, _, err := cardano.SegmentHeaderFromBytes(b)
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
  h := cardano.NewSegmentHeader(false, 0, 45)

  return testHeader(h)
}

func testTimestamped(arg string) error {
  t64, err := strconv.ParseInt(arg, 10, 64)
  if err != nil {
    return err
  }

  h := cardano.NewTimestampedSegmentHeader(t64, false, 0, 45)

  return testHeader(h)
}

func testHeader(h *cardano.SegmentHeader) error {
  h.Dump()

  b := h.ToBytes()
  fmt.Println("Integers':", b)

  hOrig, _, err := cardano.SegmentHeaderFromBytes(b)
  if err != nil {
    return err
  }

  hOrig.Dump()

  return nil
}
