package common

import (
	"bufio"
    "encoding/binary"
	"fmt"
	"net"
	"strings"
)

const DELIMITER = ","
const TERMINATOR = ";"
const BATCH_SIZE_SIZE = 4

type Message struct {
	Agency    string
	Name      string
	Surname   string
	Id        string
	Birthdate string
	Number    string
}

func (m Message) Encode() []byte {
	fields := []string{
		m.Agency,
		m.Name,
		m.Surname,
		m.Id,
		m.Birthdate,
		m.Number,
	}

	return []byte(strings.Join(fields, DELIMITER))
}

type BetSockStream struct {
	conn net.Conn
}

func BetSockConnect(address string) (*BetSockStream, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &BetSockStream{conn}, nil
}

func (s BetSockStream) PeerAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *BetSockStream) Send(msgs ...Message) error {
    encoded := make([]byte, 0)
    for _, msg := range msgs {
        encoded = append(encoded, msg.Encode()...)
    }

    batchSize := len(encoded)
    batchSizeBytes := make([]byte, 0, BATCH_SIZE_SIZE)
    binary.BigEndian.PutUint32(batchSizeBytes, uint32(batchSize))

	writer := bufio.NewWriter(s.conn)
    writer.Write(batchSizeBytes)
    writer.Write(encoded)

	err := writer.Flush()
	if err != nil {
		return fmt.Errorf("couldn't send message: %v", err)
	}

	return nil
}

func (s *BetSockStream) Close() {
	s.conn.Close()
}

type BetSockListener struct {
	listener net.Listener
}

func BetSockBind(host string, port int, backlog int) (*BetSockListener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return nil, err
	}
	return &BetSockListener{listener}, nil
}

func (l *BetSockListener) Accept() (*BetSockStream, error) {
	skt, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}
	return &BetSockStream{skt}, nil
}

func (l *BetSockListener) Close() {
	l.listener.Close()
}
