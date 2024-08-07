package pubsub

import (
	"github.com/go-redis/redis"
)

type (
	Client struct {
		*redis.Client
	}

	RedisSub struct {
		channelName string
		pubSub      *redis.PubSub
	}
)

func NewClient(addr, pass string, db int) *Client {
	var opt = redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	}
	return &Client{
		Client: redis.NewClient(&opt),
	}
}

func (c *Client) Publish(channel string, message interface{}) (int, error) {
	count, err := c.Client.Publish(channel, message).Result()
	return int(count), err
}

func (c *Client) Subscribe(channel string) (ISub, error) {
	pubSub := c.Client.Subscribe(channel)
	return &RedisSub{
		channelName: channel,
		pubSub:      pubSub,
	}, nil
}

func (c *Client) NumSub(channel string) (int, error) {
	numMap, err := c.Client.PubSubNumSub(channel).Result()
	if err != nil {
		return 0, err
	}
	if count, ok := numMap[channel]; ok {
		return int(count), nil
	}
	return 0, nil
}

func (s *RedisSub) UnSubscribe() error {
	return s.pubSub.Unsubscribe(s.channelName)
}

func (s *RedisSub) ReceiveMessage() (*Message, error) {
	message, err := s.pubSub.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	return &Message{Channel: message.Channel, Payload: message.Payload}, nil
}
