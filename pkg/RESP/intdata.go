package RESP

import "strconv"

type IntData struct {
	data int64
}

func MakeIntData(data int64) *IntData {
	return &IntData{data: data}
}

func (i *IntData) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.data, 10) + CRLF)
}
func (i *IntData) Data() int64 {
	return i.data
}
func (i *IntData) ByteData() []byte {
	return []byte(strconv.FormatInt(i.data, 10))
}
