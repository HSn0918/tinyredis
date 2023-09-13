package RESP

type ErrorData struct {
	data string
}

func MakeErrorData(data string) *ErrorData {
	return &ErrorData{data: data}
}
func (e *ErrorData) ToBytes() []byte {
	return []byte("-" + e.data + CRLF)
}

func (e *ErrorData) Error() string {
	return e.data
}
func (e *ErrorData) ByteData() []byte {
	return []byte(e.data)
}
