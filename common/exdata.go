package common

import (
	"encoding/json"
	"reflect"
)

type JsonData struct {
	Data interface{}
}

type ExData map[string]*JsonData

func (this ExData) Set(key string, value interface{}) {
	if data, ok := this[key]; ok {
		data.Data = value
	} else {
		this[key] = &JsonData{Data: value}
	}
}

func (this *JsonData) MarshalJSON() ([]byte, error) {
	switch data := this.Data.(type) {
	case []byte:
		return data, nil
	default:
		return json.Marshal(data)
	}
}

func (this *JsonData) UnmarshalJSON(data []byte) error {
	var buf []byte
	this.Data = append(buf, data...)
	return nil
}

// 把Data转换成prototype指定的类型，注意prototype为指针的指针
func (this *JsonData) To(prototype interface{}) {
	switch this.Data.(type) {
	case []byte:
		var value = reflect.New(reflect.TypeOf(prototype).Elem().Elem())
		json.Unmarshal(this.Data.([]byte), value.Interface())
		this.Data = value.Interface()
		reflect.ValueOf(prototype).Elem().Set(value)
	default:
		reflect.ValueOf(prototype).Elem().Set(reflect.ValueOf(this.Data))
	}
}
