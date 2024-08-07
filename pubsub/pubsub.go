package pubsub

// 定义了PubSub接口协议

type (
	Message struct {
		Channel string
		Payload string
	}

	// ISub 定义订阅后的操作集合
	ISub interface {
		UnSubscribe() error
		ReceiveMessage() (*Message, error)
	}

	// IPubSub 发布订阅
	IPubSub interface {
		// Publish 往频道中发布消息
		Publish(channel string, message interface{}) (int, error)
		// Subscribe 订阅一个频道
		Subscribe(channel string) (ISub, error)
		// NumSub 订阅者数量
		NumSub(channel string) (int, error)
	}
)
