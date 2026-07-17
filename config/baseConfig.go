package config

import (
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr    string              `json:"addr" redis:"addr"`
	Name    string              `json:"name" redis:"name"`
	Passwd  string              `json:"passwd" redis:"passwd"`
	TLS     bool                `json:"tls" redis:"tls"`
	Slots   []redis.ClusterSlot `json:"slots"`
	DB      int                 `json:"db" redis:"db"`
	Cluster bool                `json:"cluster" redis:"cluster"`
}

type ESConfig struct {
	Urls  []string `json:"urls" redis:"urls"`
	Sniff bool     `json:"sniff" redis:"sniff"`
}

type MysqlConfig struct {
	DSN             string `json:"dsn" redis:"dsn"`
	Debug           bool   `json:"debug" redis:"debug"`
	MaxIdleConns    int    `json:"MaxIdleConns"`
	ConnMaxLifetime int    `json:"ConnMaxLifetime"`
}

type MongoConfig struct {
	Uri string `json:"uri" redis:"uri"`
	DB  string `json:"db" redis:"db"`
}

type NetConfig struct {
	Registry string `json:"registry"`
	Nats     string `json:"nats"`
}

type BaseConfig struct {
	Release bool
	Net     NetConfig
	Mongo   MongoConfig
	Mysql   MysqlConfig
	ES      ESConfig
	Redis   RedisConfig
}

var base BaseConfig

func LoadBaseConfig(file string) error {
	return LoadToml(file, &base)
}

func LoadToml(file string, data interface{}) error {
	_, err := toml.DecodeFile(file, data)
	return err
}

func GetBase() BaseConfig {
	return base
}

func GetRegistry() string {
	return base.Net.Registry
}

func GetNats() string {
	return base.Net.Nats
}

func GetRedisConfig() RedisConfig {
	return base.Redis
}

func GetMongoConfig() MongoConfig {
	return base.Mongo
}

func GetESConfig() ESConfig {
	return base.ES
}

func GetMysqlConfig() MysqlConfig {
	return base.Mysql
}
