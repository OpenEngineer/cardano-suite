package main

import (
  "fmt"

  cardano "github.com/christianschmitz/cardano-suite"
)

func main() {
  msg := &cardano.HandshakeProposeVersions{
    map[int]int{
      0: 0,
      1: 1,
      2: 2,
    },
  }

  b := msg.ToCBOR()

  fmt.Println(b)
}
