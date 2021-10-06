package main

import (
  "fmt"

  cbor "github.com/fxamacker/cbor/v2"
  cardano "github.com/christianschmitz/cardano-suite"
)

func main() {
  h := cardano.NilHash()

  /*bList := make([]int, 8)
  for i, _ := range bList {
    bList[i] = i
  }*/

  //var untyped interface{} = bList

  var untyped interface{} = h

  b, _ := cbor.Marshal(untyped)

  fmt.Println(b)
}
