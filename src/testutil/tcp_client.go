package testutil

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	Addr string
}

func NewClient(addr string) *Client {
	return &Client{Addr: addr}
}

func (c *Client) CreateGetRequestByTCPConn(sleepTime int) (string, error) {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	time.Sleep(time.Duration(sleepTime) * time.Second)
	_, err = conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
	fmt.Println("write")
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	resp := string(buf[:n])
	return resp, nil

}
