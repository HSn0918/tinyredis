package server

import (
	"github.com/hsn/tiny-redis/pkg/RESP"
	"github.com/hsn/tiny-redis/pkg/logger"
	"github.com/hsn/tiny-redis/pkg/memdb"
	"io"
	"net"
)

type Handler struct {
	memDb   *memdb.MemDb
	aofChan chan []byte   // Channel for AOF logging
	stopCh  chan struct{} // Channel to signal shutdown
}

func NewHandler() *Handler {
	handler := &Handler{
		memDb:   memdb.NewMemDb(),
		aofChan: make(chan []byte, 100), // Buffer AOF commands
		stopCh:  make(chan struct{}),
	}
	handler.loadAOF(aofPath)
	go handler.aofLogger(aofPath)
	// 启动信号捕获协程，监听系统终止信号
	return handler
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
		// Log write commands to AOF
		if IsWriteCommand(cmd) {
			h.aofChan <- arrayData.ToBytes() // Send the command to the AOF channel
		}
	}
}
