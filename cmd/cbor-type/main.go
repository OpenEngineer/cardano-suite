// code generation for cbor serialization/deserialization
package main

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "regexp"
  "strconv"
  "strings"
)

func main() {
  if err := mainInternal(); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func autoFname(tName string) string {
  if isStruct(tName) {
    tName = tName[1:]
  }
  return tName + "_auto.go"
}

func mainInternal() error {
  args := os.Args[1:] // first is the program name, irrelevant here

  if len(args) < 1 {
    return errors.New("expected at least 1 arg (" + strings.Join(os.Args, " ") + ")")
  }

  if isStruct(args[0]) {
    return genStruct(args)
  } else {
    return genInterface(args)
  }

  return nil
}

func isList(type_ string) bool {
  return strings.HasPrefix(type_, "[]")
}

func isBuiltin(type_ string) bool {
  switch type_ {
  case "Int", "String", "IntSet", "IntList", "IntToInterfMap", "Interf", "Uint64":
    return true
  default:
    return false
  }
}

func isPtrless(type_ string) bool {
  switch type_ {
  case "ChainHash", "Hash":
    return true
  default:
    return false
  }
}

func isStruct(type_ string) bool {
  return strings.HasPrefix(type_, "*")
}

func detectPackage() (string, error) {
  pwd, err := os.Getwd()
  if err != nil {
    return "", err
  }

  files, err := ioutil.ReadDir(pwd)
  if err != nil {
    return "", err
  }

  //re := regexp.MustCompile(`^package\s*([a-z]*)`)
  re := regexp.MustCompile(`package\s([a-z]+)`)

  for _, f := range files {
    if strings.HasSuffix(f.Name(), ".go") {
      b, err := ioutil.ReadFile(f.Name())
      if err != nil {
        return "", err
      }

      str := string(b)

      res := re.FindAllStringSubmatch(str, -1)

      if len(res) == 0 || res[0][1] == "" {
        return "", errors.New("package in " + f.Name() + " not defined")
      }

      return res[0][1], nil
    }
  }

  return "", errors.New("go package not found in directory " + pwd)
}

// convenience string writer
type golangFileBuilder struct {
  fname string
  b *strings.Builder
  ind string
}

func newGolangFileBuilder(typeName string) *golangFileBuilder {
  fname := autoFname(typeName)

  bInternal := strings.Builder{}

  b := &golangFileBuilder{fname, &bInternal, ""}

  packName, err := detectPackage()
  if err != nil {
    panic(err.Error())
  }

  b.ln("package ", packName)
  b.ln()
  b.ln("import (")
  b.indent()
  b.ln("\"errors\"")
  b.ln("\"reflect\"")
  b.ln("cbor \"github.com/fxamacker/cbor/v2\"")
  b.undent()
  b.ln(")")
  b.ln()

  // avoid the compiler "imported and not used" error
  b.ln("func ", typeName, "DummyImportUsage() error {return errors.New(reflect.TypeOf(\"\").String())}")
  b.ln()

  return b
}

func (b *golangFileBuilder) ln(args ...string) {
  b.b.WriteString(b.ind)

  for _, arg := range args {
    b.b.WriteString(arg)
  }

  b.b.WriteString("\n")
}

func (b *golangFileBuilder) indent() {
  b.ind += "  "
}

func (b *golangFileBuilder) undent() {
  n := len(b.ind)

  if n < 2 {
    return 
  }

  b.ind = b.ind [0:n-2]
}

func (b *golangFileBuilder) close() error {
  return ioutil.WriteFile(b.fname, []byte(b.b.String()), 0644)
}

func genStruct(args []string) error {
  if len(args) != 2 && len(args) != 1 {
    return errors.New("expected exactly 2 arguments for converting struct")
  }

  typeName := args[0][1:]
  structFields := make([]string, 0)
  if len(args) == 2 {
    structFields = strings.Split(args[1], ",")
  }

  b := newGolangFileBuilder(typeName)

  b.ln("func ", typeName, "FromUntyped(fields interface{}) (*", typeName, ", error) {")
  b.ln("  x := &", typeName, "{}")
  b.ln("  if err := x.FromUntyped(fields); err != nil {")
  b.ln("    return nil, err")
  b.ln("  }")
  b.ln("  return x, nil")
  b.ln("}")

  b.ln()

  b.ln("func (x *", typeName, ") FromUntyped(fields_ interface{}) error {")
  b.indent()

  if len(structFields) > 0 {
    b.ln("fields, err := InterfListFromUntyped(fields_)")
    b.ln("if err != nil {")
    b.ln("  return err")
    b.ln("}")
  }

  for i, structField := range structFields {
    fieldNameType := strings.Fields(structField) 
    if len(fieldNameType) != 2 {
      return errors.New("expected only two fields in struct field spec")
    }

    fieldName := fieldNameType[0]
    fieldType := fieldNameType[1]

    iStr := strconv.Itoa(i)
    if isList(fieldType) {
      // nested lists not yet supported, only non-builtin item types supported
      itemType := fieldType[2:]
      if strings.HasPrefix(itemType, "*") {
        itemType = itemType[1:]
      }

      tmpListName := fieldName + "List"
      b.ln(tmpListName,", ok := fields[", iStr, "].([]interface{})")
      b.ln("if !ok {")
      b.ln("  return errors.New(\"expected ", fieldType, " for ", typeName, ".", fieldName, "\")")
      b.ln("}")
      b.ln(fieldName, " := make(", fieldType, ", len(", tmpListName, "))")
      b.ln("for i, item := range ", tmpListName, " {")
      b.ln("  ", fieldName, "[i], err = ", itemType, "FromUntyped(item)")
      b.ln("  if err != nil {")
      b.ln("    return err")
      b.ln("  }")
      b.ln("}")
    } else if isBuiltin(fieldType) || isPtrless(fieldType) {
      b.ln(fieldName, ", err := ", fieldType, "FromUntyped(fields[", iStr, "])")
      b.ln("if err != nil {")
      b.ln("  return errors.New(\"unexpected field type for ", typeName, ".", fieldName, ": \" + reflect.TypeOf(fields[", iStr, "]).String() + \" \" + err.Error())")
      b.ln("}")

    } else {
      if isStruct(fieldType) {
        fieldType = fieldType[1:]
      }
      b.ln(fieldName, ", err := ", fieldType, "FromUntyped(fields[", iStr, "])")
      b.ln("if err != nil {")
      b.ln("  return err")
      b.ln("}")
    }

    b.ln("x.", fieldName, " = ", fieldName)
  }

  b.ln("return nil")
  b.undent()
  b.ln("}")
  
  b.ln()

  b.ln("func (x *", typeName, ") ToUntyped() interface{} {")
  b.indent()
  b.ln("d := make([]interface{}, ", strconv.Itoa(len(structFields)), ")")
  for i, structField := range structFields {
    b.ln("{")
    b.indent()
    fieldNameType := strings.Fields(structField) 
    fieldName := fieldNameType[0]
    fieldType := fieldNameType[1]
    if isBuiltin(fieldType) {
      b.ln("var untyped interface{} = x.", fieldName)
    } else if isList(fieldType){
      b.ln("lst := make([]interface{}, len(x.", fieldName, "))")
      b.ln("for i, item := range x.", fieldName, " {")
      b.ln("  lst[i] = item.ToUntyped()")
      b.ln("}")
      b.ln("var untyped interface{} = lst")
    } else {
      b.ln("var untyped interface{} = x.", fieldName, ".ToUntyped()")
    }
    b.ln("d[", strconv.Itoa(i), "] = untyped")
    b.undent()
    b.ln("}")
  }

  b.ln("var d_ interface{} = d")
  b.ln("return d_")
  b.undent()
  b.ln("}")

  b.ln()

  // convenience CBOR functions
  b.ln("func ", typeName, "FromCBOR(b []byte) (*", typeName, ", error) {")
  b.indent()
  b.ln("d := make([]interface{}, 0)")
  b.ln("if err := cbor.Unmarshal(b, &d); err != nil {")
  b.ln("  return nil, err")
  b.ln("}")
  b.ln("var d_ interface{} = d")
  b.ln("return ", typeName, "FromUntyped(d_)")
  b.undent()
  b.ln("}")

  b.ln()

  b.ln("func (x *", typeName, ") ToCBOR() []byte {")
  b.indent()
  b.ln("d := x.ToUntyped()")
  b.ln("b, err := cbor.Marshal(d)")
  b.ln("if err != nil {")
  b.ln("  panic(err)")
  b.ln("}")
  b.ln("return b")
  b.undent()
  b.ln("}")

  return b.close()
}

func genInterface(args []string) error {
  if len(args) < 2 {
    return errors.New("expected at least 2 arguments for converting interface")
  }

  interfName := args[0]
  implementations := args[1:]

  // implementations can be structs or interfaces

  b := newGolangFileBuilder(interfName)

  b.ln("func ", interfName, "FromUntyped(fields_ interface{}) (", interfName, ", error) {")
  b.indent()
  b.ln("fields, err := InterfListFromUntyped(fields_)")
  b.ln("if err != nil {")
  b.ln("  return nil, err")
  b.ln("}")
  b.ln("t, ok := fields[0].(uint64)")
  b.ln("if !ok {")
  b.ln("  return nil, errors.New(\"first field entry isn't an int type but \" + reflect.TypeOf(fields[0]).String())")
  b.ln("}")
  b.ln("args := fields[1:]")
  b.ln("switch t {")
  for i, implName := range implementations {
    b.ln("  case ", strconv.Itoa(i), ":")
    if isStruct(implName) {
      b.ln("    return ", implName[1:], "FromUntyped(args)")
    } else {
      // must be another interface
      b.ln("    if len(args) != 1 {")
      b.ln("      return nil, errors.New(\"expected a nested list\")")
      b.ln("    }")
      b.ln("    return ", implName, "FromUntyped(args[0])")
    }
  }
  b.ln("  default:")
  b.ln("    return nil, errors.New(\"type int out of range\")")
  b.ln("}")
  b.undent()
  b.ln("}")

  b.ln()

  b.ln("func ", interfName, "ToUntyped(x_ ", interfName, ") interface{} {")
  b.indent()
  b.ln("d := make([]interface{}, 0)")
  b.ln("switch x := x_.(type) {")
  for i, implName := range implementations {
    b.ln("  case ", implName, ":")
    b.ln("    d = append(d, int(", strconv.Itoa(i), "))")

    if isStruct(implName) {
      b.ln("    tmpLst, err := InterfListFromUntyped(x.ToUntyped())")
      b.ln("    if err != nil {")
      b.ln("      panic(\"expected []interface{} from struct to untyped\")")
      b.ln("    }")
      b.ln("    d = append(d, tmpLst...)")
    } else {
      b.ln("    d = append(d, ", implName, "ToUntyped(x))")
    }
  }
  b.ln("}")
  b.ln("var d_ interface{} = d")
  b.ln("return d_")
  b.undent()
  b.ln("}")

  b.ln()

  // convenience functions
  b.ln("func ", interfName, "FromCBOR(b []byte) (", interfName, ", error) {")
  b.indent()
  b.ln("d := make([]interface{}, 0)")
  b.ln("if err := cbor.Unmarshal(b, &d); err != nil {")
  b.ln("  return nil, err")
  b.ln("}")
  b.ln("var d_ interface{} = d")
  b.ln("return ", interfName, "FromUntyped(d_)")
  b.undent()
  b.ln("}")

  b.ln()

  b.ln("func ", interfName, "ToCBOR(x ", interfName, ") []byte {")
  b.indent()
  b.ln("d := ", interfName, "ToUntyped(x)")
  b.ln("b, err := cbor.Marshal(d)")
  b.ln("if err != nil {")
  b.ln("  panic(err)")
  b.ln("}")
  b.ln("return b")
  b.undent()
  b.ln("}")

  return b.close()
}
