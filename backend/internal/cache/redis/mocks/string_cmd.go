// Package mocks provides mock implementations for Redis commands.
package mocks

import (
	"github.com/redis/go-redis/v9"
)

type StringCmd struct {
	redis.StringCmd
	val string
	err error
}

func NewStringResult(val string, err error) *StringCmd {
	return &StringCmd{val: val, err: err}
}

func (c *StringCmd) Result() (string, error) {
	return c.val, c.err
}

func (c *StringCmd) Err() error {
	return c.err
}
