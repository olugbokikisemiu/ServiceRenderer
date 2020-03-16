package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/sleekservices/ServiceRenderer/common/config"
	"github.com/sleekservices/ServiceRenderer/common/log"
)

// Client used to make requests to redis
type Client struct {
	*redis.Client

	ttl time.Duration
}

// New is a client constructor.
func New(addr string, params ...Param) *Client {
	secret := config.MustString("redis.password")
	log.Info("connecting to redis client on %s", addr)
	c := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    secret,
		DialTimeout: 5 * time.Minute,
		MaxRetries:  10,
	})

	if _, err := c.Ping().Result(); err != nil {
		log.Panic("unable to connect to redis: %s", err)
	}

	log.Info("connected to redis client")
	client := &Client{
		Client: c,
		ttl:    defaultExpirationTime,
	}

	for _, applyParam := range params {
		applyParam(client)
	}
	return client
}

func (c *Client) Ping() error {
	_, err := c.Client.Ping().Result()
	return err
}

func (c *Client) Set(key string, value interface{}) error {
	return c.Client.Set(key, value, c.ttl).Err()
}

func (c *Client) Get(key string) (interface{}, error) {
	return c.Client.Get(key).Result()
}

func (c *Client) Delete(key string) (int64, error) {
	return c.Client.Del(key).Result()
}

func (c *Client) Exists(key string) (bool, error) {
	i, err := c.Client.Exists(key).Result()
	return i >= 1, err
}

func (c *Client) SAdd(key string, value interface{}) error {
	_, err := c.Client.SAdd(key, value).Result()
	return err
}

func (c *Client) SDelete(key string) error {
	_, err := c.Delete(key)
	return err
}

func (c *Client) SRemove(key string, member interface{}) error {
	_, err := c.Client.SRem(key, member).Result()
	return err
}

func (c *Client) SMembers(key string) ([]string, error) {
	return c.Client.SMembers(key).Result()
}
