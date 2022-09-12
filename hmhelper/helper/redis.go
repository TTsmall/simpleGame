package helper

import (
	"github.com/buguang01/util"
	"github.com/garyburd/redigo/redis"
	"wawaji_pub/hmhelper/redislib"
)

var (
	redisaccess *redislib.RedisAccess
)

func Init(rdpool *redis.Pool) {
	redisaccess = new(redislib.RedisAccess)
	redisaccess.DBConobj = rdpool
}

//拿到redis的操作连接，需要主动关闭
func GetRedis() *redislib.RedisHandleModel {
	return redisaccess.GetConn()
}

func GetRedisAccess() *redislib.RedisAccess {
	return redisaccess
}

//必然释放redis连接
func UsingRedis(f func(rdmd *redislib.RedisHandleModel)) {
	rdmd := GetRedis()
	util.Using(rdmd, func() {
		f(rdmd)
	})
}
