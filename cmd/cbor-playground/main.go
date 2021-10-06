package main

import (
  "fmt"

  cbor "github.com/fxamacker/cbor/v2"

  "github.com/christianschmitz/cardano-suite/common"
)

func main() {
  h := common.NilHash()

  /*bList := make([]int, 8)
  for i, _ := range bList {
    bList[i] = i
  }*/

  //var untyped interface{} = bList

  var untyped interface{} = h

  b, _ := cbor.Marshal(untyped)

  fmt.Println(b)
}
