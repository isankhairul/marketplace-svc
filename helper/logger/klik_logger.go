package logger

import (
	"fmt"
	"github.com/go-kit/log"
	"io"
	"strings"
)

type klikLogger struct {
	w io.Writer
}

func NewKlikLogger(w io.Writer) log.Logger {
	return &klikLogger{w}
}

func (l *klikLogger) Log(keyvals ...interface{}) error {
	output, err := l.encodeKeyvals(keyvals...)
	if err != nil {
		return err
	}
	if _, e := l.w.Write([]byte(output)); e != nil {
		return e
	}
	return nil
}

func (l *klikLogger) encodeKeyvals(keyvals ...interface{}) (string, error) {
	if len(keyvals) == 0 {
		return "", nil
	}

	if len(keyvals)%2 == 1 {
		keyvals = append(keyvals, nil)
	}

	var output string
	for i := 0; i < len(keyvals); i += 2 {
		_, v := keyvals[i], keyvals[i+1]
		if _, ok := v.(string); ok && v != "" {
			output = fmt.Sprintf("%s %s", output, v)
		}
	}

	return fmt.Sprintf("%s\n", strings.TrimSpace(output)), nil
}
