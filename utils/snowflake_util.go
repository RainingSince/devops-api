package utils

import (
"github.com/sony/sonyflake"
"github.com/speps/go-hashids"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
}

func GetIntId() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

func GetUuid(prefix string) string {
	id := GetIntId()
	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + Reverse(i)
}

const Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"

func GetUuid36(prefix string) string {
	id := GetIntId()
	hd := hashids.NewData()
	hd.Alphabet = Alphabet36
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + Reverse(i)
}

