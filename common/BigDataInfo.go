package common

import (
	"database/sql/driver"
	"math/big"
)

type BigDataInfo struct {
	*big.Int
}

//有效位
const BIGINT = 100000

var BIGINTBIT = big.NewInt(BIGINT)

func NewBigIntString(s string) *BigDataInfo {
	result := new(BigDataInfo)
	result.Int = big.NewInt(0)
	result.Int.SetString(s, 10)
	return result
}

func NewBigInt(d int64) *BigDataInfo {
	result := new(BigDataInfo)
	result.Int = big.NewInt(d)
	return result
}

func NewBigIntb(b *big.Int) *BigDataInfo {
	result := new(BigDataInfo)
	result.Int = big.NewInt(0)
	result.Int.SetBits(append([]big.Word{}, b.Bits()...))
	return result
}

func (this *BigDataInfo) Decode(str string) error {
	this.Int = big.NewInt(0)
	this.Int.SetString(str, 10)
	return nil
}
func (this *BigDataInfo) MulFoalt64(f float64) {
	x := big.NewFloat(f)
	this.Int, _ = x.Mul(x, big.NewFloat(0).SetInt(this.Int)).Int(nil)
}

func (this *BigDataInfo) Scan(value interface{}) error {
	this.Int = big.NewInt(0)
	this.Int.SetString(string(value.([]byte)), 10)
	// fmt.Println("BigDataInfo:", this.Int)
	return nil
}

func (this BigDataInfo) Value() (driver.Value, error) {
	return this.MarshalText()
}

func (this *BigDataInfo) MarshalJSON() ([]byte, error) {
	if this.Int == nil {
		this.Int = big.NewInt(0)
	}
	return this.MarshalText()
}

func (this *BigDataInfo) UnmarshalJSON(b []byte) error {
	if this.Int == nil {
		this.Int = big.NewInt(0)
	}
	return this.UnmarshalText(b)
}
