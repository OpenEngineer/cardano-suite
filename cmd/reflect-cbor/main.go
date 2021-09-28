package main

import (
  "errors"
  "fmt"
  "os"
  "reflect"
  "strconv"
  "strings"

  cbor "github.com/fxamacker/cbor/v2"
)

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func dumpType(b *strings.Builder, x_ interface{}, indent string) {
  switch x := x_.(type) {
  case bool:
    if x {
      b.WriteString("true")
    } else {
      b.WriteString("false")
    }
  case uint8:
    b.WriteString(strconv.Itoa(int(x)))
  case uint64:
    b.WriteString(strconv.Itoa(int(x)))
  case int64:
    b.WriteString(strconv.Itoa(int(x)))
  case int:
    b.WriteString(strconv.Itoa(int(x)))
  case map[interface{}]interface{}:
    b.WriteString("{")
    b.WriteString("\n")
    b.WriteString(indent + "  ")

    n := len(x)
    i := 0
    for k, v := range x {
      dumpType(b, k, indent + "  ")
      b.WriteString(": ")
      dumpType(b, v, indent + "  ")

      if i < n - 1 {
        b.WriteString(", ")
      }
      i++
    }

    b.WriteString("\n")
    b.WriteString(indent)
    b.WriteString("}")
  case []interface{}:
    b.WriteString("[")
    b.WriteString("\n")
    b.WriteString(indent + "  ")

    for i, item := range x {
      dumpType(b, item, indent + "  ")

      if i < len(x) - 1 {
        b.WriteString(", ")
      }
    }

    b.WriteString("\n")
    b.WriteString(indent)
    b.WriteString("]")
  case []uint8:
    b.WriteString("[")
    b.WriteString("\n")
    b.WriteString(indent + "  ")

    for i, item := range x {
      dumpType(b, item, indent + "  ")

      if i < len(x) - 1 {
        b.WriteString(", ")
      }
    }

    b.WriteString("\n")
    b.WriteString(indent)
    b.WriteString("]")
  default:
    panic("unhandled type " + reflect.TypeOf(x).String())
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

  if len(args) < 1 {
    return errors.New("expected at least 1 arg")
  }

  b, err := parseArgsAsBytes(args)
  if err != nil {
    return err
  }

  fmt.Println("Integers:", b)

  var res interface{} = nil

  if err := cbor.Unmarshal(b, &res); err != nil {
    return err
  }

  fmt.Println("REFLECTED TYPE:")

  var sb strings.Builder

  dumpType(&sb, res, "")

  fmt.Println(sb.String())

  return nil
}
