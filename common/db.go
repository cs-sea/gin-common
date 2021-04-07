package common

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	DriverMysql = "mysql"
)

type DB struct {
	*gorm.DB
}

type Config struct {
	Driver    string
	Dsn       string
	KeepAlive int // 长连接的时间
	MaxOpen   int // 最大可以打开的数量
	MaxIdles  int // 空闲连接的最大数量
}

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

func NewDB() DB {
	gormWriter := new(Writer)

	gormLogger := NewGormLogger(*gormWriter, logger.Config{
		SlowThreshold: 0,
		Colorful:      true,
		LogLevel:      0,
	})

	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		panic(err)
	}

	return DB{db}
}

type GormLogger struct {
	Writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewGormLogger(writer Writer, config logger.Config) *GormLogger {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &GormLogger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.LogLevel = level
	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel < logger.Info {
		return
	}

	g.Printf(ctx, s, i)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel < logger.Warn {
		return
	}

	g.Printf(ctx, s, i)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel < logger.Error {
		return
	}

	g.Printf(ctx, s, i)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.LogLevel > logger.Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && g.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				g.Printf(ctx, g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(ctx, g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
			if rows == -1 {
				g.Printf(ctx, g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(ctx, g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case g.LogLevel == logger.Info:
			sql, rows := fc()
			if rows == -1 {
				g.Printf(ctx, g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(ctx, g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

type Writer struct {
	logger.Writer
}

func (w *Writer) Printf(ctx context.Context, msg string, data ...interface{}) {
	Logger(ctx).Println(
		msg,
		append([]interface{}{utils.FileWithLineNum()}, data...),
	)
}
