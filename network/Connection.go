package network

import (
  "bytes"
  "fmt"
  "net"
  "os"
  "regexp"
  "strconv"

  "github.com/christianschmitz/cardano-suite/common"
  "github.com/christianschmitz/cardano-suite/ledger"
)

const (
  HANDSHAKE     = 0
  CHAIN_SYNC    = 2
  BLOCK_FETCH   = 3
  TX_SUBMISSION = 4
  KEEP_ALIVE    = 8
  N_PROTOCOLS   = 9
)

type Connection struct {
  magic                 int
  nodeId                common.Hash
  state                 *ledger.BlockChainState
  conn                  *net.TCPConn
  incomingSegmentHeader *SegmentHeader
  incomingSegmentBuffer *bytes.Buffer 
  outgoingSegments      [][]*Segment // each protocol number gets its own segment buffer, so we can intersperse them when eventually transmitting those segments
  muxCloseIngress  chan bool   // notify the mux that the connection is closing, it is pointless to try to notify the demux like this, because is blocked during the read
  muxCloseEgress   chan bool
  muxErrorEgress   chan error  
  // demuxCloseIngress is pointless, because that thread will be blocked on conn.Read()
  demuxCloseEgress chan bool
  demuxErrorEgress chan error

  versionId     int
  versionParams interface{}

  // we use channels for passing on the messages to avoid DoS
  // messages are passed to the protocol ingresses once the buffer forms a valid CBOR message

  handshakeBuffer       *bytes.Buffer // incoming
  handshakeEgress       chan HandshakeMessage
  handshakeIngress      chan HandshakeMessage
  handshakeState        HandshakeState
  handshakeVersionEgress chan int
  handshakeCloseIngress  chan bool
  handshakeCloseEgress   chan bool
  handshakeErrorEgress   chan error

  chainSyncBuffer       *bytes.Buffer // incoming
  chainSyncEgress       chan ChainSyncMessage
  chainSyncIngress      chan ChainSyncMessage
  chainSyncState        ChainSyncState
  chainSyncCloseIngress chan bool
  chainSyncCloseEgress  chan bool
  chainSyncErrorEgress  chan error

  blockFetchBuffer       *bytes.Buffer // incoming
  blockFetchEgress       chan BlockFetchMessage
  blockFetchIngress      chan BlockFetchMessage
  blockFetchState        BlockFetchState
  blockFetchCloseIngress chan bool
  blockFetchCloseEgress  chan bool
  blockFetchErrorEgress  chan error

  txSubmissionBuffer   *bytes.Buffer // incoming
  txSubmissionEgress   chan TxSubmissionMessage
  txSubmissionIngress  chan TxSubmissionMessage

  // the keep alive protocol doesnt require a buffer for the incoming messages
  keepAliveEgress     chan KeepAliveMessage
  keepAliveIngress    chan KeepAliveMessage

  globalCloseIngress    chan bool
  globalCloseEgress     chan bool
  BlockFetchEventEgress chan *ledger.Block
  debug bool
}

// url can be an ip address or a dns address, dns addresses are resolved and the first A record that connects succesfully is used

func isDNS(host string) bool {
  re := regexp.MustCompile(`[a-zA-Z_]`)

  return re.MatchString(host)
}

func NewConnection(magic int, host string, port string, state *ledger.BlockChainState, debug bool) (*Connection, error) {
  nodeId := common.HashFromString(host)

  var (
    conn_ net.Conn
    err error
  )

  if isDNS(host) {
    addrs, err := net.LookupHost(host)

    for i, addr := range addrs {
      conn_, err = net.Dial("tcp", addr + ":" + port)
      if err != nil && i == len(addrs) - 1 {
        return nil, err
      } else if err == nil { 
        if debug {
          fmt.Println("Connected to " + addr + ":" + port)
        }
        break
      }
    }
  } else {
    conn_, err = net.Dial("tcp", host + ":" + port)
    if err != nil {
      return nil, err
    } else if err == nil && debug {
      fmt.Println("Connected to " + host + ":" + port)
    }
  }

  conn, ok := conn_.(*net.TCPConn)
  if !ok {
    panic("expected *net.TCPConn")
  }

  if err := conn.SetLinger(0); err != nil {
    return nil, err
  }

  c := &Connection{
    magic, nodeId, state, conn, nil, &bytes.Buffer{}, make([][]*Segment, N_PROTOCOLS),
    make(chan bool), make(chan bool), make(chan error), make(chan bool), make(chan error),
    0, nil,
    // handshake buffers and state
    &bytes.Buffer{}, make(chan HandshakeMessage), make(chan HandshakeMessage), HANDSHAKE_PROPOSE, make(chan int), make(chan bool), make(chan bool), make(chan error),
    &bytes.Buffer{}, make(chan ChainSyncMessage), make(chan ChainSyncMessage), CHAIN_SYNC_IDLE, make(chan bool), make(chan bool), make(chan error),
    &bytes.Buffer{}, make(chan BlockFetchMessage), make(chan BlockFetchMessage), BLOCK_FETCH_IDLE, make(chan bool), make(chan bool), make(chan error),
    &bytes.Buffer{}, make(chan TxSubmissionMessage), make(chan TxSubmissionMessage),
    make(chan KeepAliveMessage), make(chan KeepAliveMessage),
    make(chan bool), make(chan bool),
    nil,
    debug,
  }

  return c, nil
}

func (c *Connection) Version() int {
  return c.versionId
}

func (c *Connection) Listen() error {
  errCh := make(chan error)

  c.ListenNohup(errCh)

  err := <- errCh
  c.globalCloseIngress <- true
  <- c.globalCloseEgress

  return err
}

func (c *Connection) ListenNohup(errCh chan error) {
  go func() {
    c.startMux()
    c.startDemux()

    c.listenHandshake()
    c.initHandshake()

    select {
    case err := <- c.handshakeErrorEgress:
      c.closeThreads(true, false)
      errCh <- err
      <- c.globalCloseIngress
      c.globalCloseEgress <- true
      return
    case err := <- c.muxErrorEgress:
      c.closeThreads(true, false)
      errCh <- err
      <- c.globalCloseIngress
      c.globalCloseEgress <- true
      return
    case err := <- c.demuxErrorEgress:
      c.closeThreads(true, false)
      errCh <- err
      <- c.globalCloseIngress
      c.globalCloseEgress <- true
      return
    case <- c.globalCloseIngress:
      fmt.Println("[INFO] closing connection during handshake...")
      c.closeThreads(true, false)
      errCh <- nil
      c.globalCloseEgress <- true
      return
    case <- c.handshakeVersionEgress:
      fmt.Println("[INFO] handshake succeeded")
      c.handshakeCloseIngress <- true
      <- c.handshakeCloseEgress // this closes the handshake thread, but doesnt touch mux/demux
    }

    if c.debug {
      fmt.Println("[INFO] handshake resulted in protocol version", strconv.Itoa(c.versionId))
    }

    c.listenChainSync()
    c.listenBlockFetch()

    select {
    case err := <- c.chainSyncErrorEgress:
      fmt.Println("[INFO] detected chainSync error")
      c.closeThreads(false, true)
      errCh <- err
      <- c.globalCloseIngress
    case err := <- c.blockFetchErrorEgress:
      fmt.Println("[INFO] detected blockFetch error")
      c.closeThreads(false, true)
      errCh <- err
      <- c.globalCloseIngress
    case err := <- c.muxErrorEgress:
      fmt.Println("[INFO] detected mux error")
      c.closeThreads(false, true)
      errCh <- err
      <- c.globalCloseIngress
    case err := <- c.demuxErrorEgress:
      fmt.Println("[INFO] detected demux error")
      c.closeThreads(false, true)
      errCh <- err
      <- c.globalCloseIngress
    case <- c.globalCloseIngress:
      fmt.Println("[INFO] detected global close")
      c.closeThreads(false, true)
    }

    c.globalCloseEgress <- true
  }()
}

func (c *Connection) closeThreads(closeHandshake bool, closeOtherMiniProtocols bool) {
  fmt.Println("[INFO] closing some threads...")

  if c.conn == nil {
    panic("close already called before")
  }

  if closeHandshake {
    select {
    case err := <- c.handshakeErrorEgress:
      fmt.Fprintf(os.Stderr, "%s\n", err.Error())
      c.handshakeCloseIngress <- true
      <- c.handshakeCloseEgress
    case c.handshakeCloseIngress <- true: // close without pending error
      <- c.handshakeCloseEgress
    }
  } 

  if closeOtherMiniProtocols {
    select {
    case err := <- c.chainSyncErrorEgress:
      fmt.Fprintf(os.Stderr, "%s\n", err.Error())
      c.chainSyncCloseIngress <- true
      <- c.chainSyncCloseEgress
    case c.chainSyncCloseIngress <- true: // close without pending error
      <- c.chainSyncCloseEgress
    }

    if c.debug {
      fmt.Println("[INFO] chainSync thread closed")
    }

    select {
    case err := <- c.blockFetchErrorEgress:
      fmt.Fprintf(os.Stderr, "%s\n", err.Error())
      c.blockFetchCloseIngress <- true
      <- c.blockFetchCloseEgress
    case c.blockFetchCloseIngress <- true:
      <- c.blockFetchCloseEgress
    }

    if c.debug {
      fmt.Println("[INFO] blockFetch thread closed")
    }
  }

  select {
  case err := <- c.muxErrorEgress: // close with pending error
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    c.muxCloseIngress <- true
    <- c.muxCloseEgress
  case c.muxCloseIngress <- true: // close without pending error
    <- c.muxCloseEgress
  }

  if c.debug {
    fmt.Println("[INFO] mux thread closed")
  }

  c.conn.Close() // acts like the CloseIngress
  select {
  case <- c.demuxCloseEgress: // close without pending error
  case err := <- c.demuxErrorEgress: // close with pending error
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    <- c.demuxCloseEgress
  }

  if c.debug {
    fmt.Println("[INFO] demux thread closed")
    fmt.Println("[INFO] conn  thread closed")
  }

  c.conn = nil
}

func (c *Connection) Close() {
  c.globalCloseIngress <- true
  <- c.globalCloseEgress
}
