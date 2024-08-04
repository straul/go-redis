package redis

import (
	"bufio"
	"errors"
	"net"
)

func auth(conn net.Conn, password string) (string, error) {
	_, err := conn.Write([]byte("AUTH " + password + "\r\n"))
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if response[0] != '+' {
		return "", errors.New("认证失败: " + response)
	}

	return response, nil
}
