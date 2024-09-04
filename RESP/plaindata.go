package RESP

type PlainData struct {
	data string
}

func MakePlainData(data string) *PlainData {
	return &PlainData{data: data}
}

func (p *PlainData) ToBytes() []byte {
	return []byte(p.data + CRLF)
}
func (p *PlainData) Data() string {
	return p.data
}

func (p *PlainData) ByteData() []byte {
	return []byte(p.data)
}
