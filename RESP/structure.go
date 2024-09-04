package RESP

var (
	CRLF = "\r\n"
)

type RedisData interface {
	ToBytes() []byte // return RESP transfer format data
	ByteData() []byte
}
