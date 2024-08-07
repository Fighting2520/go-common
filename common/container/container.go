package container

import (
	"errors"
	"strings"
	"sync"
)

// 定义一个全局键值对存储容器
var sMap sync.Map

// CreateContainerFactory 创建一个容器工厂
func CreateContainerFactory() *containers {
	return &containers{}
}

type containers struct {
}

func (c *containers) Set(key string, value interface{}) error {
	if _, exists := c.KeyIsExists(key); !exists {
		sMap.Store(key, value)
		return nil
	}
	return errors.New("该key已在容器中存在")
}

func (c *containers) Delete(key string) {
	sMap.Delete(key)
}

func (c *containers) KeyIsExists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

func (c *containers) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(key); exists {
		return value
	}
	return nil
}

// FuzzyDelete 按照键的前缀模糊删除容器中注册的内容
func (c *containers) FuzzyDelete(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, keyPre) {
				sMap.Delete(keyName)
			}
		}
		return true
	})
}
