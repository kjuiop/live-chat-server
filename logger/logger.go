package logger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"live-chat-server/config"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

func SlogInit(cfg config.Logger) error {
	logLevel, err := slogLevelParser(cfg.Level)
	if err != nil {
		return err
	}

	var logWriter io.Writer
	if cfg.PrintStdOut {
		logWriter = io.MultiWriter(os.Stdout)
	} else {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    100, // megabytes
			MaxBackups: 10,
			MaxAge:     28,    //days
			Compress:   false, // disabled by default
		}
		logWriter = fileWriter
	}

	handler := slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = getSimplePath(2)
			}
			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return nil
}

func slogLevelParser(lvStr string) (slog.Level, error) {
	dict := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	result, ok := dict[strings.ToLower(lvStr)]
	if !ok {
		return result, fmt.Errorf("%s is not valid log level", lvStr)
	}
	return result, nil
}

func getSimplePath(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
