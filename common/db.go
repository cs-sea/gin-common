package common

import (
	"log"

	"github.com/cs-sea/gin-common/consts"
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

var logID string

func NewDB() DB {
	gormWriter := NewWriter()
	gormLogger := logger.New(gormWriter, logger.Config{})

	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		panic(err)
	}

	return DB{db}
}

type Writer struct {
	logger.Writer
	LogID string
}

func NewWriter() *Writer {
	return &Writer{
		Writer: log.New(writer("./storage/sqls/", consts.LevelInfo, 10), "", log.LstdFlags)}
}

func (w *Writer) Printf(msg string, data ...interface{}) {
	logData := map[string]interface{}{
		"log_id": logID,
	}
	data = append(data, logData)
	w.Writer.Printf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func SetLogID(logId string) {
	logID = logId
}
