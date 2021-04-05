package middlewares

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cs-sea/gin-common/common"
	"github.com/cs-sea/gin-common/consts"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (r bodyLogWriter) WriteString(s string) (n int, err error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

type LoggerRecord struct {
	startTime string
	endTime   string
	request   interface{}
	response  interface{}
	code      int
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := xid.New().String()
		common.SetLogIDField(id)
		common.SetLogID(id)

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		requestBody, _ := ioutil.ReadAll(tee)

		loggerRecord := &LoggerRecord{
			startTime: time.Now().Format(consts.DateTimeFormat),
			request:   string(requestBody),
		}

		// 重写writer 获取response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		loggerRecord.endTime = time.Now().Format(consts.DateTimeFormat)
		loggerRecord.code = blw.Status()
		loggerRecord.response = blw.body.String()

		logger := common.Logger().WithFields(map[string]interface{}{
			"startTime": loggerRecord.startTime,
			"endTime":   loggerRecord.endTime,
			"code":      loggerRecord.code,
			"request":   loggerRecord.request,
			"response":  loggerRecord.response,
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
		})

		if blw.Status() == http.StatusOK {
			logger.Infoln("request success")
		} else {
			logger.Errorln("request_error")
		}
	}
}
