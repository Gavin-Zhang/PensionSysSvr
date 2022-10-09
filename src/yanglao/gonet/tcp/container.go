package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gonet"
	"net"
	"sync"
)

const (
	MaxSocketCount uint32 = 65536
)

type tcpSocket struct {
	fd            uint32
	conn          net.Conn
	writeChan     chan []byte
	writeQuitChan chan int
	isValid       bool
	isClient      bool
	lastLen       uint32
	cacheBuffer   *bytes.Buffer
	agent         uint32
	funcName      string
	lengthBit32   bool
	readLoop      bool
}

type TcpContainer struct {
	agent             uint32
	socketSetLocker   sync.RWMutex
	socketSet         [MaxSocketCount]*tcpSocket
	listener          net.Listener
	acceptChan        chan uint32
	curConnectLocker  sync.RWMutex
	curConnect        uint32
	maxConnect        uint32
	listenLengthBit32 bool
}

func (c *TcpContainer) allocSocket() (sock *tcpSocket, err error) {
	c.socketSetLocker.Lock()
	defer func() {
		c.socketSetLocker.Unlock()
	}()

	for i := uint32(0); i < MaxSocketCount; i++ {
		if c.socketSet[i] == nil {
			sock = &tcpSocket{}
			sock.fd = i
			sock.writeChan = make(chan []byte, 1024)
			sock.writeQuitChan = make(chan int, 1)
			sock.isValid = false
			sock.isClient = false
			c.socketSet[i] = sock
			err = nil
			return
		}
	}

	return nil, errors.New("TcpConainter socket is full.")
}

func (c *TcpContainer) freeSocket(fd uint32) {
	c.socketSetLocker.Lock()
	defer func() {
		c.socketSetLocker.Unlock()
	}()

	if fd >= MaxSocketCount {
		return
	}

	sock := c.socketSet[fd]
	if sock == nil {
		return
	}

	c.socketSet[fd] = nil
}

func (c *TcpContainer) getSocket(fd uint32) (sock *tcpSocket, err error) {
	c.socketSetLocker.RLock()
	defer func() {
		c.socketSetLocker.RUnlock()
	}()

	if fd >= MaxSocketCount {
		return nil, errors.New("TcpContainer fd is out of slot.")
	}

	sock = c.socketSet[fd]
	if sock == nil {
		return nil, errors.New("TcpContainer socket is nil.")
	}

	return sock, nil
}

func (c *TcpContainer) onDisconnect(sock *tcpSocket) {
	if !sock.isValid {
		return
	}

	if !sock.isClient {
		c.curConnectLocker.Lock()
		c.curConnect--
		c.curConnectLocker.Unlock()
	}

	c.freeSocket(sock.fd)
	if sock.readLoop == true {
		sock.readLoop = false
		sock.writeQuitChan <- 0
	}

	gonet.Send(c.agent, "TcpClose", sock.fd)

	sock.fd = 0
	sock.isValid = false
	sock.isClient = false
}

func (c *TcpContainer) listenLoop() {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			break
		}

		sock, err := c.allocSocket()
		if err != nil {
			conn.Close()
			continue
		}

		c.curConnectLocker.Lock()
		if c.curConnect >= c.maxConnect {
			conn.Close()
			c.curConnectLocker.Unlock()
			continue
		}
		c.curConnect++
		c.curConnectLocker.Unlock()

		sock.conn = conn
		sock.isValid = true
		sock.lengthBit32 = c.listenLengthBit32

		// OnAccept
		gonet.Send(c.agent, "TcpAccept", sock.fd)

		go c.readLoop(sock)
		go c.writeLoop(sock)
	}
}

func (c *TcpContainer) Listen(agent uint32, addr string, maxSocket uint32, lengthBit32 bool) error {
	if c.listener != nil {
		return errors.New("TcpContainer is listened.")
	}

	laddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}

	c.listener, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	c.listenLengthBit32 = lengthBit32
	c.agent = agent
	c.maxConnect = maxSocket

	go c.listenLoop()

	return nil
}

func (c *TcpContainer) CloseListen() {
	if c.listener == nil {
		return
	}

	c.listener.Close()
}

func (c *TcpContainer) Connect(agent uint32, addr string, lengthBit32 bool) (fd uint32, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return 0, err
	}

	sock, err := c.allocSocket()
	if err != nil {
		return 0, err
	}

	sock.conn = conn
	sock.agent = agent
	sock.isValid = true
	sock.isClient = true
	sock.lengthBit32 = lengthBit32

	go c.readLoop(sock)
	go c.writeLoop(sock)

	return sock.fd, nil
}

func (c *TcpContainer) ConnectSync(agent uint32, addr string, lengthBit32 bool) (fd uint32, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return 0, err
	}

	sock, err := c.allocSocket()
	if err != nil {
		return 0, err
	}

	sock.conn = conn
	sock.agent = agent
	sock.isValid = true
	sock.isClient = true
	sock.lengthBit32 = lengthBit32

	return sock.fd, nil
}

func (c *TcpContainer) StartRWThread(fd uint32) error {
	sock, err := c.getSocket(fd)
	if err != nil {
		return err
	}

	go c.readLoop(sock)
	go c.writeLoop(sock)

	return nil
}

func (c *TcpContainer) CloseSocket(fd uint32) error {
	sock, err := c.getSocket(fd)
	if err != nil {
		return err
	}

	sock.conn.Close()
	sock.writeQuitChan <- 0
	return nil
}

func (c *TcpContainer) BindAgent(fd uint32, agent uint32) error {
	sock, err := c.getSocket(fd)
	if err != nil {
		return err
	}

	sock.agent = agent
	return nil
}

func (c *TcpContainer) Write(fd uint32, buffer []byte) error {
	sock, err := c.getSocket(fd)
	if err != nil {
		return err
	}

	sock.writeChan <- buffer
	return nil
}

func (c *TcpContainer) WriteSync(fd uint32, buffer []byte) error {
	sock, err := c.getSocket(fd)
	if err != nil {
		return err
	}

	sendBuffer := bytes.NewBuffer(nil)
	sendBytes := uint32(len(buffer))
	if sock.lengthBit32 {
		binary.Write(sendBuffer, binary.BigEndian, sendBytes)
	} else {
		binary.Write(sendBuffer, binary.BigEndian, uint16(sendBytes))
	}
	sendBuffer.Write(buffer)

	sock.conn.Write(sendBuffer.Bytes())

	return nil
}

func (c *TcpContainer) readLoop(sock *tcpSocket) {
	readBuffer := make([]byte, 1024)
	defer func() {
		c.onDisconnect(sock)
	}()

	var lastLen uint32
	cacheBuffer := bytes.NewBuffer(nil)

	for {
		readBytes, err := sock.conn.Read(readBuffer)
		if err != nil {
			break
		}

		if readBytes <= 0 {
			break
		}

		cacheBuffer.Write(readBuffer[:readBytes])

		for {
			var message []byte
			var err error

			if sock.lengthBit32 == true {
				message, err = Decode32(cacheBuffer, &lastLen)
			} else {
				message, err = Decode16(cacheBuffer, &lastLen)
			}

			if err != nil {
				// error
				sock.conn.Close()
				return
			}

			if message == nil {
				break
			}

			if sock.agent != 0 {
				gonet.Send(sock.agent, "TcpData", sock.fd, message)
			} else {
				gonet.Send(c.agent, "TcpData", sock.fd, message)
			}
		}
	}
}

func (c *TcpContainer) writeLoop(sock *tcpSocket) {
	defer func() {
		sock.readLoop = false
		c.onDisconnect(sock)
	}()

	sock.readLoop = true

	for {
		select {
		case buffer := <-sock.writeChan:
			{
				sendBuffer := bytes.NewBuffer(nil)
				sendBytes := uint32(len(buffer))
				if sock.lengthBit32 {
					binary.Write(sendBuffer, binary.BigEndian, sendBytes)
				} else {
					binary.Write(sendBuffer, binary.BigEndian, uint16(sendBytes))
				}
				sendBuffer.Write(buffer)

				bytes, err := sock.conn.Write(sendBuffer.Bytes())
				if bytes != sendBuffer.Len() {
					return
				}

				if err != nil {
					return
				}
			}
		case <-sock.writeQuitChan:
			return
		}
	}
}

func NewTcpContainer() *TcpContainer {
	container := new(TcpContainer)
	return container
}
