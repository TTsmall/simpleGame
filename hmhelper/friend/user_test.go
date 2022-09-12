package friend_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"wawaji_pub/hmhelper/friend"
)

func TestJson(t *testing.T) {
	f := new(friend.FriendMD)
	f.UserKey = "userabckey"
	m := make(map[string]StringJson)
	buf, _ := json.Marshal(f)
	k := StringJson(buf)
	m["abc"] = k
	fmt.Println(m)
	str, _ := json.Marshal(m)
	fmt.Println(string(str))
	json.Unmarshal(str, &m)
	fmt.Printf("%+v", m["abc"])
	fmt.Println("json")
	str, _ = json.Marshal(m)
	fmt.Println(string(str))

}

//json化过的字符串
type StringJson string

func (this StringJson) MarshalJSON() ([]byte, error) {
	return []byte(this), nil
}

func (this *StringJson) UnmarshalJSON(v []byte) error {
	*this = StringJson(v)
	return nil
}
