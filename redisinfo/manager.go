package redisinfo

import "context"

//异步写入redis操作
type RedisInfoManager struct {
	msgchan chan IModel
}

func (mg *RedisInfoManager) handler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(mg.msgchan)
		case msg, ok := <-mg.msgchan:
			if !ok {
				return
			}
			RedisUserInfoSet(msg)
		}
	}
}

//先不用
