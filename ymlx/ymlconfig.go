package ymlx

import (
	"errors"
	"fmt"
	"github.com/Fighting2520/go-common/container"
	"io"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

var (
	ErrInitConfigFail = errors.New("初始化配置文件错误")
)

type (
	ymlConfig struct {
		viper     *viper.Viper
		keyPrefix string
	}
)

// Clone 允许clone 一个相同功能的结构体
func (y *ymlConfig) Clone(fileName string) YMLConfiger {
	var ymlCopy = *y
	var ymlViperCopy = *(y.viper)
	ymlCopy.viper = &ymlViperCopy
	ymlCopy.viper.SetConfigName(fileName)
	if err := ymlCopy.viper.ReadInConfig(); err != nil {
		//global.ZapLog.Error(ErrInitConfigFail.Error(), zap.Error(err))
		log.Fatal(err)
	}
	return &ymlCopy
}

// Get 一个原始值
func (y *ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	}
	value := y.viper.Get(keyName)
	_ = y.cache(keyName, value)
	return value
}

// KeyIsCache 判断相关键是否已经缓存
func (y *ymlConfig) keyIsCache(keyName string) bool {
	_, exists := container.CreateContainerFactory().KeyIsExists(y.keyPrefix + keyName)
	return exists
}

// 通过键获取缓存的值
func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	return container.CreateContainerFactory().Get(y.keyPrefix + keyName)
}

// 对键值进行缓存
func (y *ymlConfig) cache(keyName string, value interface{}) error {
	return container.CreateContainerFactory().Set(y.keyPrefix+keyName, value)
}

func (y *ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	}
	value := y.viper.GetString(keyName)
	_ = y.cache(keyName, value)
	return value
}

func (y *ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	}
	value := y.viper.GetBool(keyName)
	_ = y.cache(keyName, value)
	return value
}

func (y *ymlConfig) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	}
	value := y.viper.GetInt(keyName)
	_ = y.cache(keyName, value)
	return value
}

func (y *ymlConfig) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	}
	value := y.viper.GetInt32(keyName)
	_ = y.cache(keyName, value)
	return value
}

func (y *ymlConfig) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	}
	value := y.viper.GetInt64(keyName)
	_ = y.cache(keyName, value)
	return value
}

func (y *ymlConfig) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	}
	value := y.viper.GetFloat64(keyName)
	_ = y.cache(keyName, value)
	return value
}

func (y *ymlConfig) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	}
	value := y.viper.GetDuration(keyName)
	_ = y.cache(keyName, value)
	return value
}

func (y *ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	}
	value := y.viper.GetStringSlice(keyName)
	_ = y.cache(keyName, value)
	return value
}

func CreateYamlFactoryFromReader(reader io.Reader, keyPrefix string) YMLConfiger {
	v := viper.New()
	v.SetConfigType("yml")
	if err := v.ReadConfig(reader); err != nil {
		log.Fatal(fmt.Errorf("%s, %s", ErrInitConfigFail, err))
	}
	return &ymlConfig{
		viper:     v,
		keyPrefix: keyPrefix,
	}
}

func CreateYamlFactory(configPath, keyPrefix string, filename ...string) YMLConfiger {
	yamlConfig := viper.New()
	// 配置文件所在目录
	yamlConfig.AddConfigPath(configPath)
	// 需要读取的文件名
	if len(filename) == 0 {
		yamlConfig.SetConfigName("config")
	} else {
		yamlConfig.SetConfigName(filename[0])
	}
	// 设置配置文件类型（后缀）为yml
	yamlConfig.SetConfigType("yml")
	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("%s, %s", ErrInitConfigFail, err))
	}
	return &ymlConfig{
		viper:     yamlConfig,
		keyPrefix: keyPrefix,
	}
}

// ConfigFileChangeListen 监听文件变化
func (y *ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(in fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if in.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

func (y *ymlConfig) clearCache() {
	container.CreateContainerFactory().FuzzyDelete(y.keyPrefix)
}
