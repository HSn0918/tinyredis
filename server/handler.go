package server

import (
	"github.com/hsn/tiny-redis/RESP"
	"io"
	"net"

	"github.com/hsn/tiny-redis/logger"
	"github.com/hsn/tiny-redis/memdb"
)

type Handler struct {
	memDb *memdb.MemDb
}

func NewHandler() *Handler {
	return &Handler{memDb: memdb.NewMemDb()}
}
func (h *Handler) Handle(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error(err)
		}
	}()
	ch := RESP.ParseStream(conn)
	for parsedRes := range ch {
		if parsedRes.Err != nil {
			if parsedRes.Err == io.EOF {
				logger.Info("Close connection", conn.RemoteAddr().String())
			} else {
				logger.Panic("Handle connection", conn.RemoteAddr().String(), "panic: ", parsedRes.Err.Error())
			}
			return
		}
		if parsedRes.Data == nil {
			logger.Error("empty parsedRes.Data from ", conn.RemoteAddr().String())
		}
		arrayData, ok := parsedRes.Data.(*RESP.ArrayData)
		if !ok {
			logger.Error("parsedRes.Data is not ArrayData from ", conn.RemoteAddr().String())
			continue
		}
		cmd := arrayData.ToCommand()
		res := h.memDb.ExecCommand(cmd)
		if res != nil {
			_, err := conn.Write(res.ToBytes())
			if err != nil {
				logger.Error("writer response to ", conn.RemoteAddr().String(), " error: ", err.Error())
			}
		} else {
			errData := RESP.MakeErrorData("unknown error")
			_, err := conn.Write(errData.ToBytes())
			if err != nil {
				logger.Error("writer response to ", conn.RemoteAddr().String(), " error: ", err.Error())
			}
		}
	}
}
