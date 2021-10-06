package cardano

import (
  "errors"
  "strconv"
  "strings"

  "golang.org/x/crypto/blake2b"
)

const HASH_SIZE = blake2b.Size256

type Hash [HASH_SIZE]byte

func NilHash() Hash {
  return [HASH_SIZE]byte{
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
  }
}

func HashFromBytes(data []byte) Hash {
  return blake2b.Sum256(data)
}

func HashFromString(str string) Hash {
  return HashFromBytes([]byte(str))
}

func HashFromHex(hex string) (Hash, error) {
  h := Hash{}

  if len(hex) != HASH_SIZE*2 {
    return h, errors.New("expected hex 32 bytes long")
  }

  for i := 0; i < HASH_SIZE; i++ {
    b, err := strconv.ParseUint(hex[i*2:(i+1)*2], 16, 8)
    if err != nil {
      return h, err
    }

    h[i] = byte(b)
  }

  return h, nil
}

func (h Hash) ToHex() string {
  var sb strings.Builder

  for i := 0; i < HASH_SIZE; i++ {
    sb.WriteString(strconv.FormatUint(uint64(h[i]), 16))
  }

  return sb.String()
}

func (h Hash) ToUntyped() interface{} {
  /*ints := make([]interface{}, HASH_SIZE)

  for i, b := range h {
    var b_ interface{} = uint8(b)

    ints[i] = b_
  }*/
  
  var untyped interface{} = h

  return untyped
}

func HashFromUntyped(x_ interface{}) (Hash, error) {
  x, ok := x_.([HASH_SIZE]byte)
  if !ok {
    return Hash{}, errors.New("not a [32]byte")
  }

  return Hash(x), nil

  /*for i, item_ := range x {
    item, err := IntFromUntyped(item_)
    if err != nil {
      return h, err
    }

    if item < 0 || item > 255 {
      return h, errors.New("invalid byte")
    }

    h[i] = byte(item)
  }

  return h, nil*/
}

func (h Hash) Eq(other Hash) bool {
  for i := 0; i < HASH_SIZE; i++ {
    if h[i] != other[i] {
      return false
    }
  }

  return true
}
