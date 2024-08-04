package redis

import (
	"fmt"
	"strings"
	"time"
)

func (c *Client) Set(key, value string, expiration ...time.Duration) (string, error) {
	if len(expiration) > 0 {
		exp := int(expiration[0] / time.Second)
		return c.SendCommand(fmt.Sprintf("SET %s %s EX %d", key, value, exp))
	}
	return c.SendCommand(fmt.Sprintf("SET %s %s", key, value))
}

func (c *Client) Get(key string) (string, error) {
	return c.SendCommand(fmt.Sprintf("GET %s", key))
}

func (c *Client) Del(keys ...string) (string, error) {
	return c.SendCommand(fmt.Sprintf("DEL %s", strings.Join(keys, " ")))
}

func (c *Client) Expire(key string, expiration time.Duration) (string, error) {
	exp := int(expiration / time.Second)
	return c.SendCommand(fmt.Sprintf("EXPIRE %s %d", key, exp))
}

func (c *Client) Incr(key string) (string, error) {
	return c.SendCommand(fmt.Sprintf("INCR %s", key))
}
