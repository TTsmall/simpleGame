package redislib_test

import (
	"fmt"
	"testing"

	"github.com/buguang01/Logger"
	"wawaji_pub/hmhelper/redislib"
)

func TestRedis(t *testing.T) {
	Logger.Init(0, "", Logger.LogModeFmt)
	defer Logger.LogClose()
	rd := redislib.NewRedisAccess()
	rdmd := rd.GetConn()
	// rdmd.Scan(0, "", 10)
	if result, err := rdmd.SrandMember("fruit", 2); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

}
