package main

import (
	"fmt"
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"

	"github.com/christianschmitz/cardano-suite/codec"
	"github.com/christianschmitz/cardano-suite/common"
	"github.com/christianschmitz/cardano-suite/ledger"
	"github.com/christianschmitz/cardano-suite/network"
)

func experiment1() {
	h := common.NilHash()

	var untyped interface{} = h

	b, _ := cbor.Marshal(untyped)

	fmt.Println(b)
}

// in forward sense there is not issue, because interfaces are registered
func reflectInterface(i interface{}) {
	fmt.Println(reflect.TypeOf(&i).Elem().Name())
}

func experiment2() {
	var msg network.HandshakeMessage = &network.HandshakeAcceptVersion{6, nil}

	fmt.Printf("\"%s\", \"%s\"\n", reflect.TypeOf(msg).String(), reflect.TypeOf(msg).Elem().Name())

	// a reverse lookup is possible, starting with the implementation name, and figuring out its interface,
	// implementations can thus only have one interface

	reflectInterface(msg)
}

type Struct struct {
	X   int
	Msg network.HandshakeMessage
}

func reflectStruct(i interface{}) {
	t := reflect.TypeOf(i).Elem()

	n := t.NumField()

	for i := 0; i < n; i++ {
		sf := t.Field(i)

		fmt.Println(sf.Name, sf.Type.Name())
	}
}

func experiment3() {
	str := &Struct{}

	reflectStruct(str)
}

// compare old ToUntyped to new ToUntyped
func experiment4() {
	var msg network.HandshakeMessage = network.NewHandshakeProposeVersions(network.TESTNET_MAGIC)

	b := codec.ToCBOR(msg)

	fmt.Println("#b:", b, common.ReflectCBOR(b))
}

func experiment5() {
	var msg network.HandshakeMessage = &network.HandshakeRefuse{
		&network.HandshakeRefused{6, "test"},
	}

	b := codec.ToCBOR(msg)

	fmt.Println("#b:", b, common.ReflectCBOR(b))

	var msg_ network.HandshakeMessage = nil

	if err := codec.FromCBOR(b, &msg_); err != nil {
		panic(err)
	}

	fmt.Println(msg_)

	b_ := codec.ToCBOR(msg_)

	fmt.Println("#b_:", b_, common.ReflectCBOR(b_))
}

func dumpTypeName(x interface{}) {
	v := reflect.ValueOf(x)

	fmt.Println(v.Type().Name())

	fmt.Println(reflect.Indirect(v).Type().Name())

	fmt.Println(v.Elem().Type().Name())

	fmt.Println(v.Elem().Field(0).Type().Name())

	fmt.Println(v.Elem().Field(0).Elem().Type().Name())

	fmt.Println(v.Elem().Field(0).Elem().IsNil(), v.Elem().Field(0).Elem().Type().Kind())

	fmt.Println(reflect.ValueOf(v.Elem().Field(0).Interface()).Elem().Type().Name())
}

func experiment6() {
	var msg network.HandshakeMessage = &network.HandshakeRefuse{
		&network.HandshakeRefused{6, "test"},
	}

	dumpTypeName(msg)
}

type Obj struct {
	Hello int
}

func experiment7() {
	var msg1 network.HandshakeMessage = nil

	var msg2 int

	msg2Ptr := &msg2

	v1 := reflect.ValueOf(msg1)

	v2 := reflect.ValueOf(msg2Ptr).Elem()

	var obj *Obj

	fmt.Println(v1.CanSet(), v2.CanSet(), reflect.ValueOf(&obj).CanSet())

	v2.Set(reflect.ValueOf(3))
}

func testBlockFetchMessage(msg1 network.BlockFetchMessage, msg2 network.BlockFetchMessage_) {
	b1 := network.BlockFetchMessageToCBOR(msg1)

	b2 := codec.ToCBOR(msg2)

	var msg2_ network.BlockFetchMessage_

	if err := codec.FromCBOR(b2, &msg2_); err != nil {
		panic(err)
	}

	b2 = codec.ToCBOR(msg2_)

	tStr := reflect.TypeOf(msg1).String()

	fmt.Println(b1)

	fmt.Println(b2)

	fmt.Println(common.ReflectCBOR(b1))

	fmt.Println(common.ReflectCBOR(b2))

	if len(b1) != len(b2) {
		panic(fmt.Sprintf("lengths differ for %s (expected %d, got %d)", tStr, len(b1), len(b2)))
	}

	for i, _ := range b1 {
		if b1[i] != b2[i] {
			panic(fmt.Sprintf("byte %d differs for %s", i, tStr))
		}
	}

	fmt.Printf("all good for %s\n", tStr)
}

func experiment8() {
	point := &ledger.Point{false, 2030, common.NilHash()}

	testBlockFetchMessage(
		(network.BlockFetchMessage)(&network.BlockFetchRequestRange{point, point}),
		(network.BlockFetchMessage_)(&network.BlockFetchRequestRange_{point, point}),
	)
}

func main() {
	// experiment1()

	//experiment2()

	//experiment3()

	//experiment4()

	//experiment5()

	//experiment6()

	//experiment7()

	experiment8()
}
