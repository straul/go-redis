package redis

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (c *Client) Set(ctx context.Context, key, value string, expiration ...time.Duration) (string, error) {
	if len(expiration) > 0 {
		exp := int(expiration[0] / time.Second)
		return c.SendCommand(ctx, fmt.Sprintf("SET %s %s EX %d", key, value, exp))
	}
	return c.SendCommand(ctx, fmt.Sprintf("SET %s %s", key, value))
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.SendCommand(ctx, fmt.Sprintf("GET %s", key))
}

func (c *Client) Del(ctx context.Context, keys ...string) (string, error) {
	return c.SendCommand(ctx, fmt.Sprintf("DEL %s", strings.Join(keys, " ")))
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) (string, error) {
	exp := int(expiration / time.Second)
	return c.SendCommand(ctx, fmt.Sprintf("EXPIRE %s %d", key, exp))
}

func (c *Client) Incr(ctx context.Context, key string) (string, error) {
	return c.SendCommand(ctx, fmt.Sprintf("INCR %s", key))
}
