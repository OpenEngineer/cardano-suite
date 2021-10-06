package common

import (
  "errors"
)

// Int
func IntFromUntyped(x_ interface{}) (int, error) {
  switch x := x_.(type) {
  case int:
    return x, nil
  case uint64:
    return int(x), nil
  case int64:
    return int(x), nil
  default:
    return 0, errors.New("not an integer")
  }
}

func Uint64FromUntyped(x_ interface{}) (uint64, error) {
  switch x := x_.(type) {
  case uint64:
    return x, nil
  default:
    return 0, errors.New("not a uint64")
  }
}

func InterfListFromUntyped(x_ interface{}) ([]interface{}, error) {
  x, ok := x_.([]interface{})
  if !ok {
    return nil, errors.New("expected []interf{}")
  }

  return x, nil
}

func IntListFromUntyped(x_ interface{}) ([]int, error) {
  l, ok := x_.([]interface{})
  if !ok {
    return nil, errors.New("not a list")
  }

  res := make([]int, len(l))

  for i, v := range l {
    vInt, err := IntFromUntyped(v)
    if err != nil {
      return nil, err
    }

    res[i] = vInt
  }

  return res, nil
}

// IntSet
func IntSetFromUntyped(x_ interface{}) (map[int]int, error) {
  m, ok := x_.(map[interface{}]interface{})
  if !ok {
    return nil, errors.New("not a map")
  }

  res := make(map[int]int)

  for k, v := range m {
    kInt, err := IntFromUntyped(k)
    if err != nil {
      return nil, err
    }

    vInt, err := IntFromUntyped(v)
    if err != nil {
      return nil, err
    }

    res[kInt] = vInt
  }

  return res, nil
}

func IntToInterfMapFromUntyped(x_ interface{}) (map[int]interface{}, error) {
  m, ok := x_.(map[interface{}]interface{})
  if !ok {
    return nil, errors.New("not a map")
  }

  res := make(map[int]interface{})

  for k, v := range m {
    kInt, err := IntFromUntyped(k)
    if err != nil {
      return nil, err
    }

    res[kInt] = v
  }

  return res, nil
}

func InterfFromUntyped(x_ interface{}) (interface{}, error) {
  return x_, nil
}

func StringFromUntyped(x_ interface{}) (string, error) {
  switch x := x_.(type) {
  case string:
    return x, nil
  case []byte:
    return string(x), nil
  default:
    return "", errors.New("not a string")
  }
}
