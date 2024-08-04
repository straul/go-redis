package redis

import (
	"bufio"
	"context"
	"github.com/pkg/errors"
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
func (c *Client) SendCommand(ctx context.Context, command string) (string, error) {
	conn, err := c.pool.Get()
	if err != nil {
		return "", err
	}
	defer c.pool.Put(conn)

	respChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		_, err = conn.Write([]byte(command + "\r\n"))
		if err != nil {
			errorChan <- err
			return
		}

		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			errorChan <- err
		}

		response = strings.TrimSpace(response)

		switch response[0] {
		case '+':
			respChan <- response[1:]
		case '-':
			errorChan <- errors.Errorf("redis 服务端返回错误: %s", response[1:])
		case ':':
			respChan <- response[1:]
		case '$':
			length, err := strconv.Atoi(response[1:])
			if err != nil {
				errorChan <- err
				return
			}
			if length == -1 {
				errorChan <- err
				return
			}
			data := make([]byte, length+2)
			_, err = reader.Read(data)
			if err != nil {
				errorChan <- err
				return
			}
			respChan <- string(data[:length])
		case '*':
			errorChan <- errors.New("不支持 *")
		default:
			errorChan <- errors.New("未知的返回类型")
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case response := <-respChan:
		return response, nil
	case err := <-errorChan:
		return "", err
	}
}

// Close 关闭 Redis 客户端
func (c *Client) Close() error {
	c.pool.Close()
	return nil
}
