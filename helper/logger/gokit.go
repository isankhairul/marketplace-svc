package logger

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"marketplace-svc/helper/config"
	"marketplace-svc/pkg/util"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const FileLogger = "file"

type gokitLog struct {
	logger log.Logger
}

func NewGoKitLog(cfg *config.LogConfig) Writer {
	l := gokitLog{}

	switch cfg.LogOutput {
	case FileLogger:
		l.logger = l.fileLogger(cfg.OutputFilePath)
	default:
		l.logger = l.stdLogger()
	}

	switch cfg.Level {
	case DebugLevel:
		l.logger = level.NewFilter(l.logger, level.AllowAll())
	case ErrorLevel:
		l.logger = level.NewFilter(l.logger, level.AllowError())
	case WarnLevel:
		l.logger = level.NewFilter(l.logger, level.AllowWarn())
	case InfoLevel:
		l.logger = level.NewFilter(l.logger, level.AllowInfo())
	}

	return &l
}

func (l *gokitLog) fileLogger(filePath string) log.Logger {
	logfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer func(logfile *os.File) {
		_ = logfile.Close()
	}(logfile)

	return log.NewLogfmtLogger(log.NewSyncWriter(logfile))
}

func (l *gokitLog) stdLogger() log.Logger {
	return NewKlikLogger(os.Stderr)
}

func (l *gokitLog) format(errorFormat *ErrorFormat) []interface{} {
	var keyvals []interface{}

	keyvals = append(keyvals, "ts", fmt.Sprintf("[%s]", errorFormat.DateTime.Format(util.LayoutDefault)))
	keyvals = append(keyvals, "level", strings.ToUpper(l.padding(errorFormat.Level, 5)))
	keyvals = append(keyvals, "trace-id", errorFormat.TraceID)
	keyvals = append(keyvals, "msg", errorFormat.Message)

	if errorFormat.Level != ErrorLevel {
		keyvals = append(keyvals, "caller", fmt.Sprintf("- caller=%s", errorFormat.Caller))
	}

	if errorFormat.Level == ErrorLevel {
		keyvals = append(keyvals, "stack", string(debug.Stack()))
	}

	return keyvals
}

func (l *gokitLog) padding(text string, length int) string {
	return fmt.Sprintf("%-"+strconv.Itoa(length)+"s", text)
}

func (l *gokitLog) Printf(errorFormat *ErrorFormat) {
	switch errorFormat.Level {
	case DebugLevel:
		_ = level.Debug(l.logger).Log(l.format(errorFormat)...)
	case InfoLevel:
		_ = level.Info(l.logger).Log(l.format(errorFormat)...)
	case WarnLevel:
		_ = level.Warn(l.logger).Log(l.format(errorFormat)...)
	case ErrorLevel:
		_ = level.Error(l.logger).Log(l.format(errorFormat)...)
	}
}
