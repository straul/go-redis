package redis

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Client struct {
	pool *Pool
}

// NewClient 创建 Redis 客户端
func NewClient(address, password string, poolSize int) (*Client, error) {
	pool, err := NewPool(address, password, poolSize)
	if err != nil {
		return nil, err
	}

	return &Client{pool: pool}, nil
}

// SendCommand 向 Redis 服务端发送指令并接收返回
func (c *Client) SendCommand(command string) (string, error) {
	conn, err := c.pool.Get()
	if err != nil {
		return "", err
	}
	defer c.pool.Put(conn)

	_, err = conn.Write([]byte(command + "\r\n"))
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	response = strings.TrimSpace(response)

	switch response[0] {
	case '+':
		return response[1:], nil
	case '-':
		return "", fmt.Errorf("redis 服务端返回错误: %s", response[1:])
	case ':':
		return response[1:], nil
	case '$':
		length, err := strconv.Atoi(response[1:])
		if err != nil {
			return "", err
		}
		if length == -1 {
			return "", nil
		}
		data := make([]byte, length+2)
		_, err = reader.Read(data)
		if err != nil {
			return "", err
		}
		return string(data[:length]), nil
	case '*':
		return "", errors.New("不支持 *")
	default:
		return "", errors.New("未知的返回类型")
	}
}

// Close 关闭 Redis 客户端
func (c *Client) Close() error {
	c.pool.Close()
	return nil
}
