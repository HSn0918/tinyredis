package RESP

import "strconv"

type ArrayData struct {
	data []RedisData
}

func MakeArrayData(data []RedisData) *ArrayData {
	return &ArrayData{data: data}
}
func MakeEmptyArrayData() *ArrayData {
	return &ArrayData{data: nil}
}
func (a *ArrayData) ToBytes() []byte {
	if a.data == nil {
		return []byte("*-1" + CRLF)
	}
	res := []byte("*" + strconv.Itoa(len(a.data)) + CRLF)
	for _, v := range a.data {
		res = append(res, v.ToBytes()...)
	}
	return res
}
func (a *ArrayData) Data() []RedisData {
	return a.data
}
func (a *ArrayData) ToCommand() [][]byte {
	res := make([][]byte, 0)
	for _, v := range a.data {
		res = append(res, v.ByteData())
	}
	return res
}
func (a *ArrayData) ByteData() []byte {
	res := make([]byte, 0)
	for _, v := range a.data {
		res = append(res, v.ByteData()...)
	}
	return res
}
