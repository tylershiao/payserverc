package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"reflect"
	"server/common"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liangdas/mqant/log"
	"github.com/mitchellh/mapstructure"
)

// RedisClient redis连接对象
type RedisClient struct {
	client redis.Cmdable
}

// InitRedisCilent 连接redis
func InitRedisCilent(host, name, pwd string, db int, isTls, cluster bool) (*RedisClient, error) {
	if !cluster {
		opt := &redis.Options{
			Addr:     host,
			Username: name,
			Password: pwd,
			DB:       db,
		}
		if isTls {
			opt.TLSConfig = &tls.Config{}
		}
		cli := redis.NewClient(opt)

		err := cli.Ping(context.Background()).Err()
		if err != nil {
			return nil, err
		}
		return &RedisClient{client: cli}, nil
	}
	opt := &redis.ClusterOptions{
		Addrs:    []string{host},
		Username: name,
		Password: pwd,
	}
	if isTls {
		opt.TLSConfig = &tls.Config{}
	}
	cli := redis.NewClusterClient(opt)

	err := cli.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return &RedisClient{client: cli}, nil
}

// LRange 读取一组数据，ret必须实现MarshalBinary和UnmarshalBinary方法
func (r *RedisClient) LRange(key string, ret interface{}) error {
	err := r.client.LRange(context.Background(), key, 0, -1).ScanSlice(ret)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// LPushJsonSlice 插入数组(data 必须是切片类型,元素转换成json存储)
func (r *RedisClient) LPush(key string, data interface{}) error {
	err := r.client.LPush(context.Background(), key, data).Err()
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// RPop pop出一个数据
func (r *RedisClient) RPop(key string) error {
	err := r.client.RPop(context.Background(), key).Err()
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

func (r *RedisClient) GetRedis() redis.Cmdable {
	return r.client
}

// SetJsonData 通用设置json数据
func (r *RedisClient) SetJsonData(key string, data interface{}, ex ...time.Duration) (err error) {
	ret, err := json.Marshal(data)
	if err != nil && err != redis.Nil {
		log.Error("err:%v", err)
		return
	}

	if len(ex) > 0 {
		err = r.client.Set(context.Background(), key, ret, ex[0]).Err()
	} else {
		err = r.client.Set(context.Background(), key, ret, 0).Err()
	}
	return
}

// GetJsonData 通用获取json数据
func (r *RedisClient) GetJsonData(key string, kind interface{}) (err error) {
	ret, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil && err != redis.Nil {
		log.Error("err:%v", err)
		return
	}
	err = json.Unmarshal(ret, kind)
	return
}

// SetMapData 存入一组数据
func (r *RedisClient) HSet(key string, data ...interface{}) error {
	tmp := reflect.ValueOf(data[0])
	if tmp.Kind() == reflect.Ptr {
		tmp = reflect.ValueOf(data[0]).Elem()
	}
	if tmp.Kind() == reflect.Struct {
		out := map[string]interface{}{}
		if err := mapstructure.Decode(data[0], &out); err != nil {
			log.Error("err:%v", err)
			return err
		}
		return r.GetRedis().HSet(context.Background(), key, out).Err()
	}
	return r.GetRedis().HSet(context.Background(), key, data).Err()
}

// HGetJson 获取hash，以json存储
func (r *RedisClient) HGetJson(key, field string, ret interface{}) error {
	res, err := r.GetRedis().HGet(context.Background(), key, field).Result()
	if err != nil {
		if err != redis.Nil {
			log.Error("key:%v,field:%v,err:%v", key, field, err)
		}
		return err
	}
	if err := json.Unmarshal([]byte(res), ret); err != nil {
		log.Error("key:%v,field:%v,err:%v", key, field, err)
		return err
	}
	return nil
}

// HSetJson set hash，以json存储
func (r *RedisClient) HSetJson(key, field string, ret interface{}) error {
	data, err := json.Marshal(ret)
	if err != nil {
		return err
	}
	_, err = r.GetRedis().HSet(context.Background(), key, field, data).Result()
	if err != nil {
		return err
	}
	return nil
}

// HDel 删除hash
func (r *RedisClient) HDel(key, field string) error {
	err := r.GetRedis().HDel(context.Background(), key, field).Err()
	if err != nil {
		log.Error("key:%v,field:%v,err:%v", key, field, err)
		return err
	}
	return nil
}

// HGetAll 获得某个key的所有field
func (r *RedisClient) HGetAll(key string, ret interface{}) error {
	err := r.GetRedis().HGetAll(context.Background(), key).Scan(ret)
	if err != nil {
		log.Error("HGetAll fail:%v", err)
		return err
	}
	return nil
}

// LPushJsonSlice 插入数组(data 必须是切片类型,元素转换成json存储)
func (r *RedisClient) LPushJsonSlice(key string, data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return errors.New("unknown type")
	}
	for i := 0; i < val.Len(); i++ {
		value, _ := json.Marshal(val.Index(i).Interface())
		err := r.client.LPush(context.Background(), key, string(value)).Err()
		if err != nil {
			log.Error("err:%v", err)
			return err
		}
	}
	return nil
}

// RPushJsonSlice 插入数组(data 必须是切片类型,元素转换成json存储)
func (r *RedisClient) RPushJsonSlice(key string, data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return errors.New("unknown type")
	}
	for i := 0; i < val.Len(); i++ {
		value, _ := json.Marshal(val.Index(i).Interface())
		err := r.client.RPush(context.Background(), key, string(value)).Err()
		if err != nil {
			log.Error("err:%v", err)
			return err
		}
	}
	return nil
}

// GetRankData 以score为参数倒序取set
func (r *RedisClient) GetRankData(key string) ([]redis.Z, error) {
	return r.GetRedis().ZRevRangeWithScores(context.Background(), key, 0, -1).Result()
}

// Lock redis锁
func (r *RedisClient) Lock(key string, ex ...time.Duration) bool {
	if len(ex) > 0 {
		return r.client.SetNX(context.Background(), key, 0, ex[0]).Val()
	}
	return r.client.SetNX(context.Background(), key, 1, common.RedisExpireLock).Val()
}

// UnLock redis 解锁
func (r *RedisClient) UnLock(key string, ex ...time.Duration) int64 {
	return r.client.Del(context.Background(), key).Val()
}

// WaitForLock redis等待解锁
func (r *RedisClient) WaitForLock(key string, ex ...time.Duration) {
	for {
		ok := r.Lock(key, ex...)
		if ok {
			break
		}
		// t, _ := r.GetRedis().TTL(context.Background(), key).Result()
		time.Sleep(500 * time.Millisecond)
	}
}

// SetData 通用设置key value
func (r *RedisClient) SetData(key string, data interface{}, ex ...time.Duration) error {
	if len(ex) > 0 {
		return r.GetRedis().Set(context.Background(), key, data, ex[0]).Err()
	}
	return r.GetRedis().Set(context.Background(), key, data, 0).Err()
}

func (u *RedisClient) GetInt(key string) (data int, err error) {
	data, err = u.client.Get(context.Background(), key).Int()
	if err != nil && err != redis.Nil {
		log.Error("err:%v", err)
	}
	return
}

func (u *RedisClient) GetIntWithoutErr(key string) int {
	data, err := u.client.Get(context.Background(), key).Int()
	if err != nil && err != redis.Nil {
		log.Error("err:%v", err)
		return 0
	}
	return data
}

func (u *RedisClient) GetInt64(key string) int64 {
	data, err := u.client.Get(context.Background(), key).Int64()
	if err != nil && err != redis.Nil {
		log.Error("err:%v", err)
	}
	return data
}

// GetString 通用获取string
func (r *RedisClient) GetString(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ret, err := r.GetRedis().Get(ctx, key).Result()
	return ret, err
}

func (u *RedisClient) Delkey(key string) int64 {
	res := u.client.Del(context.Background(), key)
	if err := res.Err(); err != nil {
		log.Error("err:%v", err)
	}
	count, _ := res.Result()
	return count
}

func (u *RedisClient) Exist(key string) bool {
	ret, err := u.client.Exists(context.Background(), key).Result()
	if err != nil {
		log.Error("err:%v", err)
		return false
	}
	if ret == 1 {
		return true
	}
	return false
}

// SetData 通用设置key value
func (r *RedisClient) SetNX(key string, data interface{}, ex ...time.Duration) bool {
	if len(ex) > 0 {
		return r.GetRedis().SetNX(context.Background(), key, data, ex[0]).Val()
	}
	return r.GetRedis().SetNX(context.Background(), key, data, 0).Val()
}

func (u *RedisClient) Incr(key string, value int64) (int64, error) {
	return u.client.IncrBy(context.Background(), key, value).Result()
	// if err != nil {
	// 	log.Error("err:%v", err)
	// 	return false
	// }
	// if ret == 1 {
	// 	return true
	// }
	// return false
}

// 更新玩家字段
func (u *RedisClient) HIncrBy(key, field string, inc int64) (err error) {
	if exists, _ := u.client.Exists(context.Background(), key).Result(); exists == 0 {
		return
	}
	err = u.client.HIncrBy(context.Background(), key, field, inc).Err()
	return
}

// 设置超时
func (u *RedisClient) Expire(key string, ex time.Duration) (err error) {
	err = u.client.Expire(context.Background(), key, ex).Err()
	return
}

// 获取超时
func (u *RedisClient) TTL(key string) (t time.Duration, err error) {
	return u.client.TTL(context.Background(), key).Result()
}

// HGetInt 获得某个key的制定field的int
func (r *RedisClient) HGetInt(key, field string) int {
	data, err := r.GetRedis().HGet(context.Background(), key, field).Int()
	if err != nil && err != redis.Nil {
		log.Error("HGetAll fail:%v", err)
		return 0
	}
	return data
}

func (r *RedisClient) ZAdd(key string, score int64, member interface{}) error {
	err := r.GetRedis().ZAdd(context.Background(), key, &redis.Z{Score: float64(score), Member: member}).Err()
	if err != nil {
		log.Error("ZAdd fail:%v", err)
		return err
	}
	return nil
}

func (r *RedisClient) ZAddBatch(key string, zMembers []*redis.Z) error {
	err := r.GetRedis().ZAdd(context.Background(), key, zMembers...).Err()
	if err != nil {
		log.Error("ZAdd fail:%v", err)
		return err
	}
	return nil
}

func (r *RedisClient) ZIncr(key string, score int64, member interface{}) error {
	err := r.GetRedis().ZIncr(context.Background(), key, &redis.Z{Score: float64(score), Member: member}).Err()
	if err != nil {
		log.Error("ZIncr fail:%v", err)
		return err
	}
	return nil
}

func (r *RedisClient) ZRange(key string) []redis.Z {
	ret := r.GetRedis().ZRevRangeWithScores(context.Background(), key, 0, -1)
	if ret.Err() != nil {
		log.Error("err:%v", ret.Err())
		return nil
	}
	return ret.Val()
}

func (r *RedisClient) ZRem(key string, member interface{}) error {
	err := r.GetRedis().ZRem(context.Background(), key, member).Err()
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// 按排名删除,分数高的保留
func (r *RedisClient) ZRemByRank(key string, rank int) error {
	err := r.GetRedis().ZRemRangeByRank(context.Background(), key, 0, -int64(rank+1)).Err()
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

func (r *RedisClient) ZRangeRe(key string, start, end int) []redis.Z {
	ret, err := r.GetRedis().ZRevRangeWithScores(context.Background(), key, int64(start), int64(end)).Result()
	if err != nil {
		log.Error("err:%v", err)
		return nil
	}
	return ret
}

func (r *RedisClient) ZRevRank(key string, member string) int {
	ret, err := r.GetRedis().ZRevRank(context.Background(), key, member).Result()
	if err != nil {
		log.Error("err:%v", err)
		return -1
	}
	return int(ret)
}

func (r *RedisClient) ZScore(key string, member string) int {
	ret, err := r.GetRedis().ZScore(context.Background(), key, member).Result()
	if err != nil {
		log.Error("err:%v", err)
		return -1
	}
	return int(ret)
}

func (r *RedisClient) GeyKeys(key string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ret, err := r.GetRedis().Keys(ctx, key).Result()
	if err != nil {
		log.Error("err:%v", err)
		return nil
	}
	return ret
}
