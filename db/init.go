package db

import (
	"server/config"
	edb "server/db/es"
	mdb "server/db/mysql"
	rdb "server/db/redis"

	"github.com/liangdas/mqant/log"
)

var (
	redisClient *rdb.RedisClient
	esClient    *edb.EsClient
	mysqlClient *mdb.MysqlClient
)

// InitDB 初始化db
func InitDB(dbs ...interface{}) {
	var err error
	for _, v := range dbs {
		switch v.(type) {
		case *rdb.RedisClient:
			c := config.GetRedisConfig()
			redisClient, err = rdb.InitRedisCilent(c.Addr, c.Name, c.Passwd, c.DB, c.TLS, c.Cluster)
			if err != nil {
				log.Error("err:%v", err)
				panic("connect redis fail")
			}
		case *edb.EsClient:
			c := config.GetESConfig()
			esClient, err = edb.InitEsClient(c.Urls, c.Sniff)
			if err != nil {
				log.Error("err:%v", err)
				panic("connect es fail")
			}
		case *mdb.MysqlClient:
			c := config.GetMysqlConfig()
			mysqlClient, err = mdb.InitMysqlClient(c.DSN, c.Debug)
			if err != nil {
				log.Error("err:%v", err)
				panic("connect mysql fail")
			}
		}
	}
}

func Redis() *rdb.RedisClient {
	return redisClient
}

func Mysql() *mdb.MysqlClient {
	return mysqlClient
}

func ES() *edb.EsClient {
	return esClient
}

func SetRedis(r *rdb.RedisClient) {
	redisClient = r
}

func SetDB(r *mdb.MysqlClient) {
	mysqlClient = r
}

func SetES(r *edb.EsClient) {
	esClient = r
}
