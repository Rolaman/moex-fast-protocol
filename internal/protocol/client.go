package protocol

import (
	"net"
	"strconv"

	"golang.org/x/net/ipv4"
)

// According to FAST specification (1.2.5)
const maxUDPPackageSize = 1500

type Client struct {
	options *Options

	conn         *ipv4.PacketConn
	netInterface *net.Interface
	group        *net.UDPAddr
	source       *net.UDPAddr
	name         string

	IsConnected bool
}

func NewClient(options *Options) (*Client, error) {
	c := Client{
		options: options,
	}

	en0, err := net.InterfaceByName("ppp0")
	if err != nil {
		return nil, err
	}
	c.netInterface = en0

	packetConn, err := net.ListenPacket("udp4", net.JoinHostPort("0.0.0.0", strconv.Itoa(options.Port)))
	if err != nil {
		return nil, err
	}
	c.conn = ipv4.NewPacketConn(packetConn)
	c.group = &net.UDPAddr{IP: net.ParseIP(options.Group)}
	c.source = &net.UDPAddr{IP: net.ParseIP(options.Source)}

	return &c, nil
}

func (c *Client) Connect() error {
	if err := c.conn.JoinSourceSpecificGroup(c.netInterface, c.group, c.source); err != nil {
		return err
	}
	c.IsConnected = true
	return nil
}

func (c *Client) Disconnect() error {
	if err := c.conn.LeaveSourceSpecificGroup(c.netInterface, c.group, c.source); err != nil {
		return err
	}
	c.IsConnected = false
	return nil
}

func (c *Client) ReadNext() ([]byte, error) {
	buf := make([]byte, maxUDPPackageSize)
	if _, _, _, err := c.conn.ReadFrom(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Client) ReadNextN() ([][]byte, error) {
	buffer := make([]ipv4.Message, 5)
	for i := 0; i < 5; i++ {
		buffer[i] = ipv4.Message{
			Buffers: [][]byte{
				make([]byte, maxUDPPackageSize),
			},
		}
	}
	if _, err := c.conn.ReadBatch(buffer, 0); err != nil {
		print(err)
		return nil, err
	}
	response := make([][]byte, 1)
	for i, msg := range buffer {
		if buffer[i].Buffers[0][0] != 0 && buffer[i].Buffers[0][1] != 0 {
			if i == 0 {
				response[i] = msg.Buffers[0]
			} else {
				response = append(response, msg.Buffers[0])
			}
		} else {
			if i == 0 {
				return make([][]byte, 0), nil
			}
		}
	}
	return response, nil
}

func (c *Client) String() string {
	if c.name != "" {
		return c.name
	}
	return c.options.Group + ":" + strconv.Itoa(c.options.Port)
}
