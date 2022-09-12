package hmrank

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/buguang01/util"
	"github.com/buguang01/util/threads"
	"wawaji/hmhelper/helper"
	"wawaji/hmhelper/redislib"
)

type RankManage struct {
	rankconflist map[string]*RankConfig //排行榜配置
	gamename     string                 //游戏前缀，用来在写入redis的时候使用
	redismd      *redislib.RedisAccess
	msgchan      chan IRankWriteModel
	thg          *threads.ThreadGo
}

func (md *RankManage) Init() error {
	return nil
}
func (md *RankManage) Run() {
	fmt.Println("hmrank.RankManage Run")
	md.thg.Go(md.autoTask)
	return
}
func (md *RankManage) Stop() {
	md.thg.CloseWait()
	return
}

//新建排行榜管理器
func NewRankManage(redismd *redislib.RedisAccess, gamename string, rankcf map[string]*RankConfig) *RankManage {
	result := new(RankManage)
	result.redismd = redismd
	result.gamename = gamename
	result.rankconflist = rankcf
	result.msgchan = make(chan IRankWriteModel, 1024)
	result.thg = threads.NewThreadGo()
	return result
}

func (mg *RankManage) autoTask(ctx context.Context) {
	now := time.Now()
	//周期排行榜，删明天要用的列表
	next := helper.GetDate(now).Add(time.Hour*24 - 5*time.Minute)
	// fmt.Println("autoTask", next.Sub(now).Seconds())
	tk := time.NewTimer(next.Sub(now))
	// fmt.Println("autoTask", now.String(), next.String())
	//有效时间排行榜，删昨天的列表
	next2 := helper.GetDate(now).Add(time.Hour * 24)
	tk2 := time.NewTimer(next2.Sub(now))
	rand.Seed(int64(now.Nanosecond()))
	randtk := time.NewTimer(time.Duration(rand.Int63n(int64(time.Second))) + time.Minute)
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case msg := <-mg.msgchan:
			{
				rdkeyli := msg.GetRedisKey()
				ukey := msg.GetUserKey()
				rdscore := msg.GetRedisScore() //有时间的值
				score := msg.GetScore()        //没有时间的值
				rdmd := mg.redismd.GetConn()
				util.Using(rdmd, func() {
					for _, key := range rdkeyli {
						rdkey := fmt.Sprint(mg.gamename, "_", key)
						rdmd.Zadd(rdkey, rdscore, ukey)
						rdinfokey := fmt.Sprint(mg.gamename, "_info_", key)
						rdmd.Hmset(rdinfokey, ukey, score)
					}
				})
			}
		case td := <-tk.C:
			{
				// fmt.Println("autoTask周期排行榜，删明天要用的列表:", td.String())
				rdmd := mg.redismd.GetConn()
				util.Using(rdmd, func() {
					tmday := td.AddDate(0, 0, 1)
					for _, cf := range mg.rankconflist {
						if cf.RkCycle != -1 {
							//明天与今天是同一个周期
							if cf.GetDayCycleNum(td) == cf.GetDayCycleNum(tmday) {
								continue
							}
							//周期排行榜，删明天要用的列表
							cday := cf.GetDayCycleNum(tmday)
							// cday := tmday / cf.RkCycle
							cday = cday % cf.RkValidity
							rdkey := fmt.Sprint(mg.gamename, "_", cf.GetRedisKeyByDay(cday))
							rdmd.Del(rdkey)
						}
					}
				})

				now = time.Now()
				next = next.AddDate(0, 0, 1)
				tk.Reset(next.Sub(now))
				// fmt.Println("autoTask22周期排行榜，删明天要用的列表", next.Sub(now).Seconds())
			}
		case td2 := <-tk2.C:
			{
				// fmt.Println("autoTask有效时间排行榜，删昨天的列表:", td2.String())
				rdmd := mg.redismd.GetConn() //没有时间的值
				util.Using(rdmd, func() {
					// tmday := int(td2.Add(-time.Hour*24).Unix() / 86400)
					for _, cf := range mg.rankconflist {
						tmday := cf.GetDayNum(td2.AddDate(0, 0, -1))
						if cf.RkCycle == -1 && cf.RkValidity != -1 {
							//有效时间排行榜，删昨天的列表
							cday := tmday % cf.RkValidity
							rdkey := fmt.Sprint(mg.gamename, "_", cf.GetRedisKeyByDay(cday))
							rdmd.Del(rdkey)
						}
					}
				})

				now = time.Now()
				next2 = next2.AddDate(0, 0, 1)
				tk2.Reset(next2.Sub(now))
			}
		case <-randtk.C:
			{
				//每分钟刷新随机种子
				now = time.Now()
				rand.Seed(int64(now.Nanosecond()))
				randtk.Reset(time.Duration(rand.Int63n(int64(time.Second))) + time.Minute)
			}
			// default:
			// 	{
			// 		time.Sleep(1 * time.Second)
			// 		now := time.Now()
			// 		fmt.Println("default", next.Sub(now).Seconds(), next2.Sub(now).Seconds())
			// 	}
		}
	}
}

//写入chan
func (mg *RankManage) SendRankModel(msg IRankWriteModel) {
	mg.msgchan <- msg
}
