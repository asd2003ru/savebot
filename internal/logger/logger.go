package logger

//package logger is simple wrapper for zerolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	stdlog "log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	TIME_FORMAT = "2006-01-02T15:04:05.000Z"
)

// Dumper interface
type IDumper interface {
	// Add line to multiline log
	Add(msg string, args ...any) *DumpLog
	// Add object as JSON to multiline log
	AddJson(v any) *DumpLog
	// Commit multiline log as error
	CommitAsError()
	//Commit multiline log as info
	CommitAsInfo()
	// Commit multiline log as debug
	CommitAsDebug()
	// Clear multiline log
	Clear() *DumpLog
}

// ILogger interface
type ILogger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(err error, msg string, args ...any)
	Fatal(err error, msg string, args ...any)
	Panic(err error, msg string, args ...any)
	Trace(msg string, args ...any)
	// WithFields adds custom fields to the logger
	//
	//	log.WithFields(log.Fields{"foo": "bar"}).Info("Hello, World!")
	WithFields(Fields) ILogger
	// StdLogger returns the standard logger
	StdLogger() *stdlog.Logger
	GetDumper() IDumper
}

type Fields map[string]any

type Zlog struct {
	zlog *zerolog.Logger
}

type DumpLog struct {
	msgs strings.Builder
	zl   *Zlog
}

// New multiline logger
func (zl *Zlog) GetDumper() IDumper {
	return &DumpLog{msgs: strings.Builder{}, zl: zl}
}

// Add to multiline log
func (dmp *DumpLog) Add(msg string, args ...any) *DumpLog {
	dmp.msgs.WriteString(fmt.Sprintf("%s\n", fmt.Sprintf(msg, args...)))
	return dmp
}

// Add to multiline log as JSON
func (dmp *DumpLog) AddJson(v any) *DumpLog {
	b, _ := json.MarshalIndent(v, "", " ")
	dmp.msgs.WriteString("\n============ JSON ============\n")
	dmp.msgs.WriteString(string(b))
	dmp.msgs.WriteString("\n============ ---- ============\n")
	return dmp
}

// Commit as error
func (dmp *DumpLog) CommitAsError() {
	dmp.zl.Error(nil, "%s", dmp.msgs.String())
}

// Commit as info
func (dmp *DumpLog) CommitAsInfo() {
	dmp.zl.Info("%s", dmp.msgs.String())
}

// Commit as debug
func (dmp *DumpLog) CommitAsDebug() {
	dmp.zl.Debug("%s", dmp.msgs.String())
}

// Clear multiline log
func (dmp *DumpLog) Clear() *DumpLog {
	dmp.msgs.Reset()
	return dmp
}

// Write to log (for StdLogger)
func (zl *Zlog) Write(p []byte) (n int, err error) {
	p = bytes.TrimSpace(p)
	zl.Info("%s", string(p))
	return len(p), nil
}

// Get standard logger
func (zl *Zlog) StdLogger() *stdlog.Logger {

	return stdlog.New(zl, "STDLOG: ", 0)
}

// Fatal log.
func (zl *Zlog) Fatal(err error, msg string, args ...any) {

	zl.zlog.Fatal().Err(err).Str("src", sourcePoint()).Msgf(msg, args...)
}

// Info log.
func (zl *Zlog) Info(msg string, args ...any) {
	zl.zlog.Info().Msgf(msg, args...)
}

// Panic log.
func (zl *Zlog) Panic(err error, msg string, args ...any) {

	zl.zlog.Panic().Err(err).Str("src", sourcePoint()).Msgf(msg, args...)
}

// Trace log.
func (zl *Zlog) Trace(msg string, args ...any) {

	zl.zlog.Trace().Str("src", sourcePoint()).Msgf(msg, args...)
}

// WithFields implements ILogger.
func (zl *Zlog) WithFields(fields Fields) ILogger {
	zev := zl.zlog.With()
	if len(fields) > 0 {
		for k, v := range fields {
			if e, ok := v.(error); ok {
				zev = zev.AnErr(k, e)
				continue
			}
			zev = zev.Interface(k, v)

		}
	}
	zlog := zev.Logger()
	return &Zlog{
		zlog: &zlog,
	}

}

// Error implements ILogger.
func (zl *Zlog) Error(err error, msg string, args ...any) {

	zl.zlog.Error().Err(err).
		Str("src", sourcePoint()).
		Msgf(msg, args...)

}

// Debug implements ILogger.
func (zl *Zlog) Debug(msg string, args ...any) {

	zl.zlog.Debug().Str("src", sourcePoint()).Msgf(msg, args...)

}
func (zl *Zlog) Warn(msg string, args ...any) {

	zl.zlog.Warn().Str("src", sourcePoint()).Msgf(msg, args...)

}

func sourcePoint() string {
	_, file, line, _ := runtime.Caller(2)
	s := 0
	short := ""
	for n := len(file) - 1; n > 0; n-- {
		if file[n] == '/' || file[n] == '\\' {
			s += 1
			short = file[n+1:]
		}
		if s > 3 {
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

// Type for log. Json for service and text for local debug
//
// Use Env for override (CLOG_TYPE=0 or 1)
type LogType uint8

const (
	// Json default log
	JSONType LogType = iota
	// TextType for local debug
	TextType
)

// Level defines log levels (zlog).
type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
	// Values less than TraceLevel are handled as numbers.
)

// New logger instance.
//
//	logtype: LogType (JSONType or TextType) (use Env CLOG_TYPE=text for override)
//
//	loglevel: Log level (use Env CLOG_LEVEL=debug,info, etc for override)
func NewLogger(logtype LogType, loglevel Level) ILogger {

	// Change log type if Env is set
	if lt := os.Getenv("CLOG_TYPE"); lt != "" && strings.ToLower(lt) == "text" {
		logtype = 1
	}

	// Set log level from Env if set
	zerolog.SetGlobalLevel(zerolog.Level(loglevel))
	if ll := os.Getenv("CLOG_LEVEL"); ll != "" {
		lv, err := zerolog.ParseLevel(ll)
		if err == nil {
			zerolog.SetGlobalLevel(lv)
		}
	}

	zerolog.TimeFieldFormat = TIME_FORMAT
	var _zlog zerolog.Logger
	if logtype == TextType {
		_zlog = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: TIME_FORMAT})
	} else {
		_zlog = log.Logger
	}

	return &Zlog{
		zlog: &_zlog,
	}

}
