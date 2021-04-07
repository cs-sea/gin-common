package common

import (
	"context"
	"path"
	"time"

	"github.com/cs-sea/gin-common/consts"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const FileSuffix = ".log"
const LogPath = "./storage/logs"

var globalLogger *logrus.Logger

const (
	LogIDKey = "log_id"
)

func init() {
	globalLogger = logrus.New()
	// 分割文件配置
	NewSimpleLogger(globalLogger, LogPath, 10)
}

func NewSimpleLogger(log *logrus.Logger, logPath string, save uint) {

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer(logPath, "debug", save), // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer(logPath, "info", save),
		logrus.WarnLevel:  writer(logPath, "warn", save),
		logrus.ErrorLevel: writer(logPath, "error", save),
		logrus.FatalLevel: writer(logPath, "fatal", save),
		logrus.PanicLevel: writer(logPath, "panic", save),
	}, &logrus.JSONFormatter{TimestampFormat: consts.DateTimeFormat})

	log.AddHook(lfHook)
}

func writer(logPath string, level string, save uint) *rotatelogs.RotateLogs {
	logFullPath := path.Join(logPath, level)

	logier, err := rotatelogs.New(
		logFullPath+"-%Y%m%d%H"+FileSuffix,
		rotatelogs.WithLinkName(logFullPath), // 生成软链，指向最新日志文件
		//rotatelogs.WithRotationCount(save),       // 文件最大保存份数
		rotatelogs.WithMaxAge(time.Hour*24*7),  // 最大天数 这俩个只能用一个
		rotatelogs.WithRotationTime(time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		globalLogger.Errorln(err)
	}
	return logier
}

func Logger(ctx context.Context) *logrus.Entry {
	entry := logrus.NewEntry(globalLogger)
	return entry.WithField(LogIDKey, ctx.Value(LogIDKey))
}
