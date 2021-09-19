package protocol

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

type TcpClient struct {
	options *TcpOptions
	conn    net.Conn
	reader  *bufio.Reader
}

func NewTcpClient(options *TcpOptions) *TcpClient {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", options.Host, options.Port))
	buf := bufio.NewReader(conn)
	return &TcpClient{
		options: options,
		conn:    conn,
		reader:  buf,
	}
}

func (c *TcpClient) Send(b []byte) []byte {
	_, err := c.conn.Write(b)
	if err != nil {
		panic(err)
	}
	resp, err := ioutil.ReadAll(c.reader)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *TcpClient) SendAndReadUntilMessage(until func([]byte) bool, b []byte) [][]byte {
	_, err := c.conn.Write(b)
	if err != nil {
		panic(err)
	}
	resps := make([][]byte, 0)
	resp, err := ioutil.ReadAll(c.reader)
	resps = append(resps, resp)
	if until(resp) {
		resp, err := ioutil.ReadAll(c.reader)
		if err != nil {
			log.Println(err)
			return resps
		}
		resps = append(resps, resp)
	}
	return resps
}
