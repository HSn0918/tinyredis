package RESP

type StringData struct {
	data string
}

func MakeStringData(data string) *StringData {
	return &StringData{
		data: data,
	}
}
func (r *StringData) ToBytes() []byte {
	return []byte("+" + r.data + CRLF)
}
func (r *StringData) Data() string {
	return r.data
}
func (r *StringData) ByteData() []byte {
	return []byte(r.data)
}
