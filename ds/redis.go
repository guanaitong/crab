package ds

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/guanaitong/crab/gconf"
	"github.com/guanaitong/crab/system"
	"github.com/guanaitong/crab/util/format"
	"strings"
)

const defaultRedisConfigKey = "redis-config.json"

type RedisType int

const (
	RedisStandalone RedisType = iota
	RedisSentinel
)

// 单机模式配置
type StandaloneConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	NodeHost string `json:"nodeHost"`
	NodePort int    `json:"nodePort"`
}

// 哨兵模式配置
type SentinelConfig struct {
	Master string `json:"master"`
	Nodes  string `json:"nodes"`
}

type RedisConfig struct {
	Type              RedisType        `json:"type"`
	Standalone        StandaloneConfig `json:"standalone"`
	Sentinel          SentinelConfig   `json:"sentinel"`
	Password          string           `json:"password"`
	EncryptedPassword string           `json:"encryptedPassword"`
	Db                int              `json:"db"`
}

func (redisConfig *RedisConfig) NewClient() *redis.Client {
	var pwd = decrypt(redisConfig.EncryptedPassword)
	if pwd == "" {
		pwd = redisConfig.Password
	}

	if redisConfig.Type == RedisStandalone {
		var (
			host, port = redisConfig.Standalone.Host, redisConfig.Standalone.Port
		)
		//预发环境使用nodeHost和nodePort
		if system.GetWorkEnv() == "prepare" && system.GetWorkIdc() == "sh" && redisConfig.Standalone.NodeHost != "" {
			host, port = redisConfig.Standalone.NodeHost, redisConfig.Standalone.NodePort
		}
		opt := &redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: pwd,
			DB:       redisConfig.Db,
		}
		return redis.NewClient(opt)
	} else if redisConfig.Type == RedisSentinel {
		fOpt := &redis.FailoverOptions{
			MasterName:    redisConfig.Sentinel.Master,
			SentinelAddrs: strings.Split(redisConfig.Sentinel.Nodes, ","),
			Password:      pwd,
			DB:            redisConfig.Db,
		}
		return redis.NewFailoverClient(fOpt)
	} else {
		panic("unsupported type")
	}
}

func GetRedisConfig(key string) *RedisConfig {
	if key == "" {
		panic(errors.New("redis config is null"))
	}

	redisConfig := new(RedisConfig)
	configValue := gconf.GetCurrentConfigCollection().GetConfig(key)
	err := format.AsJson(configValue, redisConfig)
	if err != nil {
		panic(err.Error())
	}
	return redisConfig
}

func GetDefaultRedisConfig() *RedisConfig {
	return GetRedisConfig(defaultDataSourceKey)
}
