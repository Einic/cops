/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: logger_type.go
 * @Version: 1.0.0
 * @Date: 2024/1/21 20:42
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	globalLogger *zap.SugaredLogger
	logMutex     sync.Mutex
	logChan      chan LogMsg
)

const (
	CdmAsstRunLogFile = "./logs/cops-run.log"
)

// zapLogger implements the Logger interface
type zapLogger struct {
	sugarLogger *zap.SugaredLogger
}

type LogMsg struct {
	zapcore.Entry
	Level   zapcore.Level
	Message string
	Fields  []zap.Field
}
