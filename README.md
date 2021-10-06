# Cardano in golang
My goals for this repository:
* library of cardano protocol functions and utilities (but not for a fully functional stake pool)
* cmd-line utilities
* utility to compile golang code into plutus core
* wallet with ui for some very specific smart contracts
* disprove  the 'formal methods' narrative of IOHK by creating a non-staking cardano node which is simpler to analyse than IOHK's haskell code

## Types of smart contracts
* Staked Testament or "Last will", UTXO in this contract can be used any time for any payment
* Stable coins with staked reserves
  * minted stable coins have a time-stamp
  * only longer held stable coins are completely redeemable
  * each reserve-provider initiates a different contract (different staking addres), the input datum determines the contract parameters, the contracts are listable by their script hash
* Maker type exchange contracts for on-chain assets (wallet needs smart functionality for cost-efficient Takers)
  * no need for DEX orderbook
* Insurance, loans, identity, password manager, 

## Comparison with cardano-node haskell implementation
* building cardano-node from scratch drags in about 100K lines of code
* I'm sure we can de better, without so many dependencies

## Bootstrapping the network

# Networking 

## p2p network
The p2p manager/governor maintains a list of three kinds of peers
* cold peers (this could be the entire list of relay nodes discovered first in the bootstrap topology file, and later in block-chain itself)
* warm/pingable peers
* hot peers

It has functions to demote/promote each kind of peer to each preceding/subsequent kind of peer.
It should maintain a target number of each type of peer, and the set of warm peers should be diverse in terms of hop distance and geographic locations.

The lists should also be updated with a minimum frequency (i.e. a target churn frequency)

For each peer some reputation is maintained, in case connections fail.

Incoming connections aren't possible because this is a client library and we only consume data.
So we are abusing the cardano network protocol for a client-server relationship only

Questions:
* which bootstrap node should we connect to?
  * see: https://developers.cardano.org/docs/get-started/running-cardano
  * basically specifies the following node: relays-new.cardano-testnet.iohkdev.io:3001
* how do we get a list of other nodes?
  * longer topology file
  * addresses are registered in block-chain with pool certificates?
    * what if one stake pool has several relay nodes -> should be visible in same pool certificate (also some dns records can contain multiple peers/ip addresses, but for now just assume one)
    * so first we must sync significant parts of the block-chain before being able to discover peers

## Single p2p connection
A list of long-running TCP connections between peers

Each connection has a multiplexer

The (de)multiplexer sees packets and splits them in data segments

Could a (de)multiplexer be seen as a golang channel?

**A first cmd-line could be a tool to connect with a single node**

Different mini-protocols are run through the multiplexer:
* chain-sync
* block-fetch
* tx-submission
* keep-alive


### Handshake mini-protocol

Reverse engineering of ouroboros-network-framework/src/Ouroboros/Network/Protocol/Handshake/Codec.hs

Note: builtin (de)serialization of Golang wrt. any codec is a huge plus for Golang

* Client proposes versions by sending a list to the Server
* The server can also simultaneously send a version proposal
* The server either accepts a version, or refuses
* Then the client either refuses or confirms the version

How would you implement a TCP-ping-pong protocol in golang? The client type has a certain state and evolves upon each received message, optional returning a message to be sent to the server.

So a protocol is initiated into the connection. A protocol has a certain number and will receive all messages with that number.

The protocol defines 3 types of messages:
* MsgProposeVersions
* MsgRefuse
* MsgAcceptVersion

Assume that all CBOR functions in Codec.hs simply append the parts with the `<>` operator

#### MsgProposeVersions
[
  byte // msg subtype (0), 
  { // versions
    version-number::int(32 or 64?) /* 0, 1 or 2 according to Test.hs */ -> CBOR.Term
  }
]

Feels like map is just used as a set, and type map value is left unspecified?

#### MsgAcceptVersion
[
  byte, // msg subtype (1),
  version-number,
  param, // type and content irrelevant?
]

#### MsgRefuse
[
  byte // msg subtype (2),
  [ // VersionMismatch
    0 (byte of error subtype), [list-of-actually-accepted-version-numbers],
  ] or [ // HandshakeDecodeError
    1 (byte of error subtype), version-number, error-string,
  ] or [ // Refused
    2 (byte of error subtype), version-number, refuse-reason-string,
  ]
]

#### Handshake type deserialization
Assuming the byte-type pattern is reused throughout the ouroboros protocol, there should be a way to dynamically specify the 

The general interface HandshakeMessage, with three implementations:
* HandshakeProposeVersions
* HandshakeAcceptVersion
* HandshakeRefuse

The following suffixes, ToCBOR and FromCBOR, serialize and deserialize respectively

The tool to generate these functions requires the complete types of course. Preferably in golang source code:

```
type HandshakeMessage interface {
  <methods>
}

// different code is needed in case of child interface or child struct
//go:generate cbor-type HandshakeMessage *HandshakeProposeVersions *HandshakeAcceptVersion HandshakeRefuse
=> 
func HandshakeMessageFromUntyped(fields []interface{}) HandshakeMessage {
  <unpack one level of nested interface slices>
  t, ok := fields[0].byte
  <handle ok>

  args := fields[1:]

  switch t {
    case 0:
      return HandshakeProposeVersionsFromUnTyped(args)
      
  }
}
// auto-generate convenience method
func HandshakeMessageFromCBOR(d []byte) HandshakeMessage {
}

func HandshakeMessageToUntyped(x HandshakeMessage) []interface{} {
  // first entry in output list is the type byte

  // remaining entries are the fields in case of a struct
}


//go:generate cbor-type *HandshakeProposeVersions "versions map[int]int"
//go:generate cbor-type *HandshakeAcceptVersion "version int, param int"
//go:generate cbor-type  HandshakeRefuse *HandshakeVersionMismatch *HandshakeDecodeError *HandshakeRefused

// the cbor-struct takes a comma separated list and recognizes builtin types
//go:generate cbor-struct HandshakeVersionMismatch "versions []int"
=>  // not a member method so the function can be called generically
func HandshakeVersionMismatchFromUntyped(args []interface{}) *HandshakeVersionMismatch {
  x := &HandshakeVersionMismatch{}

  // for builtin type
  versions, ok := args[0].([]int)
  <handle ok>

  x.versions = versions

  return x
}

// convenience method
func (x *HandshakeVersionMismatch) ToCBOR() []byte {
  return <convert CBOR to byte slice>
}

//go:generate cbor-type *HandshakeDecodeError "version int, reason string"
//go:generate cbor-type *HandshakeRefused     "version int, error string"

```

## Epoch vs slot vs block
An epoch contains multiple slots, but not every slot contains a block

## Blockchain state
Syncing Daedelus after a week offline takes about 50 minutes. This is slow and annoying.
Instead we can just pick some nodes we trust to act as servers, and use them to 'vote' on the state of the blockchain. This approach completely ignores the PoS consensus mechanism of Cardano, but is probably good enough for a wallet tracking the state of the blockchain. Only block hashing is required, which can done as a scanning operation as we move through the blockchain.

The connections are initiated with a number of nodes trusted by the user.

Each node starts with a clean slate, and if errors occur they are added to a score.

Nodes are identified by a 32 bit hash of their domain-names. Each node gets an equal weight.

The local header state is kept in memory completely

Each dissenting block header is also tracked.

Blocks themselves are not stored.

Not everything needs to be stored in memory. Only the last 10000 or so blocks, so each header can store extended meta data (i.e. a list of which nodes confirmed)

Use GOB to serialize/deserialize the in-memory block state.
Forks should also be tracked. Each fork has a pointer to a location on the main chain, the main chain is the one with the most votes for the first block.

Forks are removed once their prev block is no longer in memory, or once their first block doesn't have any more votes. What about a fork of a fork?

Using BlockIds is much faster than hash lookup or

Block queries are performed using the main chain.

The number of hot nodes changes all the time. So only count the number of hot votes.

Another thread polls the main chain to download the block data.

The main chain can only be shortened once its first block has been downloaded and verified.

Shortening is triggered when a certain length is exceeded.

## TODO:
* CBOR serializer/deserializer code-generator: done
* Handshake tester: done
* Connection type implementing all the Node-to-Node mini protocols
* Remove duplicate code in mini protocols using golang generics
