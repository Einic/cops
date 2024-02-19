/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: logger.go
 * @Version: 1.0.0
 * @Date: 2024/1/26 15:32
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package zaplog

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/wxnacy/wgo/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

// The Logger interface defines the basic operations of the logger
type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Panic(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	Close()
}

// InitLogger initializes the logger
func InitLogger(callerSkip int) Logger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()

	logChan = make(chan LogMsg, 1000)
	go processLogMessages()

	// Use callerSkip to skip the appropriate number of frames
	logger := zap.New(
		zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel),
		zap.AddCaller(),
		// +1 Used to account for extra skips in this function
		zap.AddCallerSkip(callerSkip+1),
	)
	globalLogger = logger.Sugar()

	return &zapLogger{sugarLogger: globalLogger}
}

func (zl *zapLogger) log(entry zapcore.Entry, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Create an interface slice to store fields
	var interfaceFields []interface{}
	for _, f := range fields {
		interfaceFields = append(interfaceFields, f)
	}

	// Use sugarLogger directly to record logs
	switch entry.Level {
	case zapcore.DebugLevel:
		zl.sugarLogger.Debugw(entry.Message, interfaceFields...)
	case zapcore.InfoLevel:
		zl.sugarLogger.Infow(entry.Message, interfaceFields...)
	case zapcore.WarnLevel:
		zl.sugarLogger.Warnw(entry.Message, interfaceFields...)
	case zapcore.ErrorLevel:
		zl.sugarLogger.Errorw(entry.Message, interfaceFields...)
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		zl.sugarLogger.Panicw(entry.Message, interfaceFields...)
	case zapcore.FatalLevel:
		zl.sugarLogger.Fatalw(entry.Message, interfaceFields...)
	}
}

func (zl *zapLogger) Debug(message string, fields ...zap.Field) {
	zl.log(zapcore.Entry{Level: zapcore.DebugLevel, Message: message}, fields...)
}

func (zl *zapLogger) Info(message string, fields ...zap.Field) {
	zl.log(zapcore.Entry{Level: zapcore.InfoLevel, Message: message}, fields...)
}

func (zl *zapLogger) Warn(message string, fields ...zap.Field) {
	zl.log(zapcore.Entry{Level: zapcore.WarnLevel, Message: message}, fields...)
}

func (zl *zapLogger) Error(message string, fields ...zap.Field) {
	zl.log(zapcore.Entry{Level: zapcore.ErrorLevel, Message: message}, fields...)
}

func (zl *zapLogger) Panic(message string, fields ...zap.Field) {
	zl.log(zapcore.Entry{Level: zapcore.PanicLevel, Message: message}, fields...)
}

func (zl *zapLogger) Fatal(message string, fields ...zap.Field) {
	zl.log(zapcore.Entry{Level: zapcore.FatalLevel, Message: message}, fields...)
}

func (zl *zapLogger) Close() {
	_ = zl.sugarLogger.Sync()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:          "ts",
			LevelKey:         "level",
			NameKey:          "logger",
			CallerKey:        "caller_line",
			FunctionKey:      zapcore.OmitKey,
			MessageKey:       "msg",
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      cEncodeLevel,
			EncodeTime:       cEncodeTime,
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     cEncodeCaller,
			ConsoleSeparator: " ",
		})
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   CdmAsstRunLogFile,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     1,
		Compress:   true,
	}
	ws := io.MultiWriter(lumberJackLogger, os.Stdout)
	return zapcore.AddSync(ws)
}

func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level.CapitalString() {
	case "DEBUG":
		enc.AppendString(color.Cyan(level.CapitalString()))
	case "INFO":
		enc.AppendString(color.Green(level.CapitalString()))
	case "WARN":
		enc.AppendString(color.Yellow(level.CapitalString()))
	case "ERROR", "DPANIC", "PANIC":
		enc.AppendString(color.Red(level.CapitalString()))
	case "FATAL":
		enc.AppendString(color.BgRed(level.CapitalString()))
	default:
		enc.AppendString(level.CapitalString())
	}
}

func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}

func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%-17s", caller.TrimmedPath()))
}

func processLogMessages() {
	for logMsg := range logChan {
		switch logMsg.Entry.Level {
		case zapcore.DebugLevel:
			globalLogger.Debugw(logMsg.Entry.Message, toInterfaceSlice(logMsg.Fields)...)
		case zapcore.InfoLevel:
			globalLogger.Infow(logMsg.Entry.Message, toInterfaceSlice(logMsg.Fields)...)
		case zapcore.WarnLevel:
			globalLogger.Warnw(logMsg.Entry.Message, toInterfaceSlice(logMsg.Fields)...)
		case zapcore.ErrorLevel:
			globalLogger.Errorw(logMsg.Entry.Message, toInterfaceSlice(logMsg.Fields)...)
		case zapcore.DPanicLevel, zapcore.PanicLevel:
			globalLogger.Panicw(logMsg.Entry.Message, toInterfaceSlice(logMsg.Fields)...)
		case zapcore.FatalLevel:
			globalLogger.Fatalw(logMsg.Entry.Message, toInterfaceSlice(logMsg.Fields)...)
		}
	}
}

// Convert zap.Fields to []interface{}
func toInterfaceSlice(fields []zap.Field) []interface{} {
	interfaceFields := make([]interface{}, len(fields))
	for i, f := range fields {
		interfaceFields[i] = f
	}
	return interfaceFields
}
