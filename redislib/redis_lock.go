package redislib

//redis的事务模式
func (rd *RedisHandleModel) RedisWatch(keys []interface{}, f func() error) error {
	rd.Do("WATCH", keys...)
	defer rd.Do("UNWATCH")
	return f()
}
