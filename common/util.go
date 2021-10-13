package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	cbor "github.com/fxamacker/cbor/v2"
)

// prefer a
func Min(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

// prefer a
func Max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func detach(fn func() error, onError func()) {
	go func() {
		if err := fn(); err != nil {
			fmt.Fprintf(os.Stderr, "Connection error: %s\n", err.Error())
			fmt.Fprintf(os.Stderr, "Closing connection\n")
			onError()
		}
	}()
}

func AppendToCBORBuffer(buf *bytes.Buffer, b []byte, maxTotalSize int) ([]interface{}, error) {
	if _, err := buf.Write(b); err != nil {
		return nil, err
	}

	if buf.Len() > maxTotalSize {
		return nil, errors.New("packet buffer overflow")
	}

	// make a copy of the buffer, and use that to decode
	bufCopy := &bytes.Buffer{}

	bufCopy.Write(buf.Bytes())

	dec := cbor.NewDecoder(bufCopy)

	data := make([]interface{}, 0)

	for {
		var item interface{}

		if err := dec.Decode(&item); err != nil {
			if err == io.EOF {
				// advance the original buffer
				dummy := make([]byte, dec.NumBytesRead())
				if _, err := buf.Read(dummy); err != nil {
					panic(err)
				}

				return data, nil
			} else {
				return nil, err
			}
		} else {
			data = append(data, item)
		}
	}
}

func reflectUntyped(b *strings.Builder, x_ interface{}, indent string) {
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
	case [HASH_SIZE]byte:
		b.WriteString("Hash[...]\n")
	case map[interface{}]interface{}:
		b.WriteString("{")
		b.WriteString("\n")
		b.WriteString(indent + "  ")

		n := len(x)
		i := 0
		for k, v := range x {
			reflectUntyped(b, k, indent+"  ")
			b.WriteString(": ")
			reflectUntyped(b, v, indent+"  ")

			if i < n-1 {
				b.WriteString(", ")
			}
			i++
		}

		b.WriteString("\n")
		b.WriteString(indent)
		b.WriteString("}")
	case []interface{}:
		b.WriteString("[")
		b.WriteString(strconv.Itoa(len(x)))
		b.WriteString("]")
		b.WriteString("[")
		b.WriteString("\n")
		b.WriteString(indent + "  ")

		for i, item := range x {
			reflectUntyped(b, item, indent+"  ")

			if i < len(x)-1 {
				b.WriteString(", ")
			}
		}

		b.WriteString("\n")
		b.WriteString(indent)
		b.WriteString("]")
	case []uint8:
		if len(x) == HASH_SIZE {
			for _, i := range []int{0, 1, 30, 31} {
				b.WriteString(strconv.FormatUint(uint64(x[i]), 16))
				if i == 1 {
					b.WriteString("..")
				}
			}
		} else if len(x) == 64 {
			for _, i := range []int{0, 1, 2, 61, 62, 63} {
				b.WriteString(strconv.FormatUint(uint64(x[i]), 16))
				if i == 2 {
					b.WriteString("..")
				}
			}
		} else {
			b.WriteString("Uint8[")
			b.WriteString(strconv.Itoa(len(x)))
			b.WriteString("][")
			b.WriteString("\n")
			b.WriteString(indent + "  ")

			for i, item := range x {
				reflectUntyped(b, item, indent+"  ")

				if i < len(x)-1 {
					b.WriteString(", ")
				}
			}

			b.WriteString("\n")
			b.WriteString(indent)
			b.WriteString("]")
		}
	case string:
		b.WriteString("\"")
		b.WriteString(x)
		b.WriteString("\"")
	case cbor.Tag:
		b.WriteString("cbor.Tag{")
		b.WriteString(strconv.Itoa(int(x.Number)))
		b.WriteString(",")

		reflectUntyped(b, x.Content, indent+"  ")

		b.WriteString("\n")
		b.WriteString(indent)
		b.WriteString("}")
	default:
		panic("unhandled type " + reflect.TypeOf(x).String())
	}
}

func ReflectUntyped(x interface{}) string {
	var sb strings.Builder

	reflectUntyped(&sb, x, "")

	return sb.String()
}

// note: b might contain multiple cbor messages
func ReflectCBOR(b []byte) string {
	buf := &bytes.Buffer{}

	buf.Write(b)

	dec := cbor.NewDecoder(buf)

	var (
		err error
		sb  strings.Builder
	)

	nBefDecoding := buf.Len()

	for err == nil {
		var res interface{} = nil

		err = dec.Decode(&res)
		if err == nil {
			if err_ := cbor.Unmarshal(b, &res); err_ != nil {
				panic(err_)
			}

			reflectUntyped(&sb, res, "")
		} else {
			nOk := dec.NumBytesRead()

			for i := 0; i < nBefDecoding-buf.Len()-nOk; i++ {
				fmt.Println("#unreading byte ", i)
				if err := buf.UnreadByte(); err != nil {
					panic(err)
				}
			}

			fmt.Println("#decode encountered an error: ", err.Error(), buf.Len(), nOk)
			return sb.String()
		}
	}

	return sb.String()
}

// used in WrappedBlock and WrappedBlockHeader
func UnwrapCBOR(x_ interface{}) (interface{}, error) {
	x, ok := x_.(cbor.Tag)
	if !ok {
		return nil, errors.New("not a cbor.Tag")
	}

	if x.Number != 24 {
		return nil, errors.New("not a cbor.Tag{24}")
	}

	data, ok := x.Content.([]byte)
	if !ok {
		return nil, errors.New("cbor.Tag content not a byte list")
	}

	var res interface{} = nil

	if err := cbor.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func WrapCBOR(x interface{}) (interface{}, error) {
	b, err := cbor.Marshal(x)
	if err != nil {
		return nil, err
	}

	var untyped interface{} = cbor.Tag{24, b}

	return untyped, nil
}
