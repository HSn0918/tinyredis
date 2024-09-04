package logger

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	PANIC
)

func Debug(v ...any) {
	if logcfg.Level > DEBUG {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(DEBUG)
	logger.Println(v)
}
func Info(v ...any) {
	if logcfg.Level > INFO {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(INFO)
	logger.Println(v)
}

func Warning(v ...any) {
	if logcfg.Level > WARNING {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...any) {
	if logcfg.Level > ERROR {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(ERROR)
	logger.Println(v)
}

func Panic(v ...any) {
	if logcfg.Level > PANIC {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(PANIC)
	logger.Println(v)
}
