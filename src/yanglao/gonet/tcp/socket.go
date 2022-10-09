package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gonet"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type TcpAdapter interface {
	onConnClose(conn *tcpConn)
	onRead(bytes int)
	onWrite(bytes int)
}

type tcpConn struct {
	adapter       TcpAdapter
	fd            uint32
	socket        net.Conn
	agent         uint32
	writeChan     chan []byte
	writeQuitChan chan int
	connected     bool
	decoder       *PackageDecoder
}

func (c *tcpConn) attach(agent uint32) {
	c.agent = agent
}

func (c *tcpConn) Process(message []byte) {
	gonet.Send(c.agent, "TcpData", message)
}

func (c *tcpConn) DecodeError() {
	fmt.Printf("tcpConn DecodeError")
	c.adapter.onConnClose(c)
}

////////////////////////////////////////
// TcpServer
type TcpServer struct {
	listener          *net.TCPListener
	agent             uint32
	maxClient         uint32
	clientCount       uint32
	clientCountLocker sync.RWMutex
	slot              []*tcpConn

	netMonitor bool
	read       uint64
	write      uint64
}

func (s *TcpServer) loop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		tcpConn := &tcpConn{s, 0, conn, 0, make(chan []byte, 16), make(chan int), false, NewPackageDecoder()}
		s.onAccept(tcpConn)
	}
}

func (s *TcpServer) onAccept(conn *tcpConn) {
	s.clientCountLocker.Lock()
	defer s.clientCountLocker.Unlock()

	if s.clientCount == s.maxClient {
		conn.socket.Close()
		return
	}

	s.clientCount++

	var fd uint32 = 0
	for fd = 0; fd < s.maxClient; fd++ {
		if s.slot[fd] == nil {
			s.slot[fd] = conn
			break
		}
	}
	conn.connected = true
	conn.fd = fd
	gonet.Send(s.agent, "TcpServerAccept", fd)

	go tcpConnRead(s, conn)
	go tcpConnWrite(s, conn)
}

func (s *TcpServer) onConnClose(conn *tcpConn) {
	if !conn.connected {
		return
	}

	s.clientCountLocker.Lock()
	defer s.clientCountLocker.Unlock()
	s.clientCount--

	s.slot[conn.fd] = nil
	conn.connected = false

	conn.socket.Close()
	//conn.writeQuitChan <- 0

	gonet.Send(s.agent, "TcpServerClose", conn.fd)
}

func (s *TcpServer) onRead(bytes int) {
	if s.netMonitor {
		s.read += uint64(bytes)
	}
}

func (s *TcpServer) onWrite(bytes int) {
	if s.netMonitor {
		s.write += uint64(bytes)
	}
}

func (s *TcpServer) Attach(fd uint32, agent uint32) {
	conn := s.slot[fd]
	if conn == nil {
		fmt.Println("Attach: conn is nil.")
		return
	}
	conn.attach(agent)
}

func (s *TcpServer) Write(fd uint32, buffer []byte) {
	conn := s.slot[fd]
	if (conn != nil) && conn.connected {
		sendBuffer := bytes.NewBuffer(nil)
		sendBytes := uint16(len(buffer))
		binary.Write(sendBuffer, binary.BigEndian, sendBytes)
		sendBuffer.Write(buffer)
		conn.writeChan <- sendBuffer.Bytes()
	}
}

func (s *TcpServer) Close(fd uint32) {
	conn := s.slot[fd]
	if conn != nil {
		//s.onConnClose(conn)
		conn.writeQuitChan <- 0
	}
}

func (s *TcpServer) GetClientCount() uint32 {
	return s.clientCount
}

func (s *TcpServer) LocalAddr(fd uint32) (addr net.Addr) {
	conn := s.slot[fd]
	if conn != nil {
		addr = conn.socket.LocalAddr()
		return
	}

	return
}

func (s *TcpServer) RemoteAddr(fd uint32) (addr net.Addr) {
	conn := s.slot[fd]
	if conn != nil {
		addr = conn.socket.RemoteAddr()
		return
	}

	return
}

func (p *TcpServer) SetNetMonitor(isOpen bool) {
	if p.netMonitor == isOpen {
		return
	}

	p.netMonitor = isOpen
	p.read = 0
	p.write = 0
	if p.netMonitor {
		time.AfterFunc(time.Second, p.NetMonitor)
	}
}

func (p *TcpServer) NetMonitor() {
	log.Println(fmt.Sprintf("AcceptCount:%d|ReadBytes:%v|WirteBytes:%v", p.clientCount, p.read, p.write))
	p.read, p.write = 0, 0
	if p.netMonitor {
		time.AfterFunc(time.Second, p.NetMonitor)
	}
}

func (p *TcpServer) ListenClose() bool {
	p.listener.Close()

	for _, conn := range p.slot {
		if conn != nil {
			conn.socket.Close()
		}
	}
	return true
}

func NewTcpServer(agent uint32, addr string, maxClient uint32) *TcpServer {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	server := &TcpServer{}
	server.listener = listener
	server.agent = agent
	server.maxClient = maxClient
	server.slot = make([]*tcpConn, maxClient)
	server.netMonitor = false

	go server.loop()

	return server
}

////////////////////////////////////////
// TcpClient
type TcpClient struct {
	conn  *tcpConn
	agent uint32
}

func (c *TcpClient) LocalAddr() net.Addr {
	return c.conn.socket.LocalAddr()
}

func (c *TcpClient) RemoteAddr() net.Addr {
	return c.conn.socket.RemoteAddr()
}

func (c *TcpClient) Write(buffer []byte) {
	if (c.conn != nil) && c.conn.connected {
		sendBuffer := bytes.NewBuffer(nil)
		sendBytes := uint16(len(buffer))
		binary.Write(sendBuffer, binary.BigEndian, sendBytes)
		sendBuffer.Write(buffer)
		c.conn.writeChan <- sendBuffer.Bytes()
	}
}

func (c *TcpClient) Close(fd uint32) {
	c.onConnClose(c.conn)
}

func (c *TcpClient) onConnClose(conn *tcpConn) {
	if !conn.connected {
		return
	}
	conn.connected = false
	conn.writeQuitChan <- 0

	gonet.Send(c.agent, "TcpClientClose")
}

func (c *TcpClient) onRead(bytes int) {

}

func (s *TcpClient) onWrite(bytes int) {

}

func NewTcpClient(agent uint32, addr string) *TcpClient {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
        return nil
    }

	client := &TcpClient{}
	client.agent = agent
	client.conn = &tcpConn{client, 0, conn, agent, make(chan []byte, 16), make(chan int), false, NewPackageDecoder()}
	client.conn.connected = true

	go tcpConnRead(client, client.conn)
	go tcpConnWrite(client, client.conn)

	return client
}

func checkError(err error) {
	if err != nil {
		log.Printf("Tcp Error: %s", err.Error())
		os.Exit(1)
	}
}

////////////////////////////////////////
// share read/write
func tcpConnRead(adapter TcpAdapter, conn *tcpConn) {
	readBuffer := make([]byte, 1024)
	defer func() {
		adapter.onConnClose(conn)
	}()

	for {
		readBytes, err := conn.socket.Read(readBuffer)
		adapter.onRead(readBytes)
		if err != nil {
			break
		}

		if readBytes == 0 {
			break
		}

		err = conn.decoder.Decode(readBuffer[:readBytes], conn)
		if err != nil {
			break
		}
	}
}

func tcpConnWrite(adapter TcpAdapter, conn *tcpConn) {
	// defer adapter.onConnClose(conn)
	defer func() {
		adapter.onConnClose(conn)
	}()

	for {
		select {
		case buffer := <-conn.writeChan:
			bytes, err := conn.socket.Write(buffer)
			adapter.onWrite(bytes)
			if bytes != len(buffer) {
				return
			}
			if err != nil {
				return
			}
		case <-conn.writeQuitChan:
			return
		}
	}
}
