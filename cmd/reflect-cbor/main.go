package main

import (
  "errors"
  "fmt"
  "os"
  "regexp"
  "strconv"
  "strings"

  "github.com/christianschmitz/cardano-suite/common"
)

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
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

func parseArgsAsInts(args []string) ([]byte, error) {
  b := make([]byte, len(args))

  for i, arg := range args {
    v, err := strconv.ParseInt(arg, 10, 64)
    if err != nil {
      return nil, errors.New("failed to parse byte " + strconv.Itoa(i))
    }

    b[i] = byte(v)
  }

  return b, nil
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

func stripCommas(args []string) []string {
  res := make([]string, len(args))

  for i, arg := range args {
    res[i] = strings.Trim(arg, ",")
  }

  return res
}

func mainInternal() error {
  args := os.Args[1:]

  if len(args) < 1 {
    return errors.New("expected at least 1 arg")
  }

  var (
    b []byte
    err error
  )
  
  args = stripCommas(args)

  if allInts(args) {
    b, err = parseArgsAsInts(args)
  } else {
    b, err = parseArgsAsBytes(args)
    if err != nil {
      return err
    }
  }

  fmt.Println(fmt.Sprintf("#Raw integers[%d]:", len(b)), b)

  str, err := common.ReflectCBOR(b)

  fmt.Println(str)

  if err != nil {
    return err
  }


  return nil
}
