package logger

import (
	"context"
	"fmt"
	"marketplace-svc/app/model/request"
	"marketplace-svc/pkg/util"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	validation "github.com/itgelo/ozzo-validation/v4"
)

const (
	// InfoLevel is the default logging priority.
	InfoLevel = "info"

	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel = "warn"

	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = "error"

	//DebugLevel Anything else, i.e. too verbose to be included in INFO level.
	DebugLevel = "debug"

	DefaultDepthCaller = 3
)

type Writer interface {
	Printf(errorFormat *ErrorFormat)
}

type ErrorFormat struct {
	DateTime   time.Time
	Level      string
	TraceID    string
	Message    string
	Caller     string
	StackTrace string
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(err error)
	Debug(err error)
	Handle(ctx context.Context, err error)
	WithContext(ctx context.Context) Logger
}

type logger struct {
	Writer      Writer
	Context     context.Context
	DepthCaller int
}

func NewLogger(writer Writer) Logger {
	return &logger{
		Writer:      writer,
		DepthCaller: DefaultDepthCaller,
	}
}

func (l *logger) WithContext(ctx context.Context) Logger {
	lg := *l
	lg.Context = ctx
	return &lg
}

func (l *logger) Info(msg string) {
	l.Writer.Printf(l.Log(msg, InfoLevel))
}

func (l *logger) Warn(msg string) {
	l.Writer.Printf(l.Log(msg, WarnLevel))
}

func (l *logger) Error(err error) {
	l.Writer.Printf(l.Log(err.Error(), ErrorLevel))
}

func (l *logger) Debug(err error) {
	l.Writer.Printf(l.Log(err.Error(), DebugLevel))
}

func (l *logger) Log(msg string, level string) *ErrorFormat {
	var traceId string
	if l.Context != nil {
		if ctxTraceID, ok := l.Context.Value(TraceIDContextKey).(string); ok {
			traceId = ctxTraceID
		}
	}

	return &ErrorFormat{
		DateTime: util.TimeNow(),
		Level:    level,
		TraceID:  traceId,
		Message:  msg,
		Caller:   l.caller(l.DepthCaller),
	}
}

func (l *logger) caller(depth int) string {
	_, file, line, _ := runtime.Caller(depth)
	idx := strings.LastIndexByte(file, '/')
	return file[idx+1:] + ":" + strconv.Itoa(line)
}

func DefaultRawLogFormat(errorFormat *ErrorFormat) string {
	var stackTrace string
	if errorFormat.Level == DebugLevel {
		stackTrace = fmt.Sprintf("- %s", string(debug.Stack()))
	}

	var caller string
	if errorFormat.Level != DebugLevel {
		caller = fmt.Sprintf(" - caller=%s", errorFormat.Caller)
	}

	var traceFormat string
	if errorFormat.TraceID != "" {
		traceFormat = fmt.Sprintf("trace-id=%s ", errorFormat.TraceID)
	}

	return fmt.Sprintf("[%s] %s %s%s%s %s",
		errorFormat.DateTime.Format(util.LayoutDefault),
		strings.ToUpper(errorFormat.Level),
		traceFormat,
		errorFormat.Message,
		caller,
		stackTrace,
	)
}

// Handle - implement for ServerErrorHandler
func (l *logger) Handle(_ context.Context, err error) {
	switch err.(type) {
	default:
		l.Error(err)
	case validation.Errors, request.MalformedRequest:
		l.Info(err.Error())
	}
}
