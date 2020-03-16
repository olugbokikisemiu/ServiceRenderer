package redis

import "time"

const (
	defaultExpirationTime = time.Hour
)

// Param is an optional param for redis client.
type Param func(*Client)

// WithTTL used to set keys expiration time.
func WithTTL(t time.Duration) Param {
	return func(c *Client) {
		c.ttl = t
	}
}
