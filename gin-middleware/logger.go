package gin_middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xmchz/go-common/util"
	"time"
)

func InitAccessLogger(name string, relativePath string) (*log.Logger, chan map[string]interface{}) {
	ch := make(chan map[string]interface{}, 100)
	logger := util.NewLogger(name, relativePath, &log.TextFormatter{})
	return logger, ch
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func AccessLog(logger *log.Logger, ch chan map[string]interface{}) gin.HandlerFunc {
	go func() {
		for accessLog := range ch {
			logger.WithFields(accessLog).Info()
		}
	}()

	return func(c *gin.Context) {

		startTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()
		accessLogMap := make(map[string]interface{})
		accessLogMap["request_time"] = startTime
		accessLogMap["response_time"] = endTime
		accessLogMap["cost_time"] = fmt.Sprintf("%vms", endTime-startTime)
		accessLogMap["request_method"] = c.Request.Method
		accessLogMap["response_code"] = c.Writer.Status()
		accessLogMap["request_uri"] = c.Request.RequestURI
		//accessLogMap["request_proto"] = c.Request.Proto
		//accessLogMap["request_ua"] = c.Request.UserAgent()
		//accessLogMap["request_referer"] = c.Request.Referer()

		//body,_ := c.Request.GetBody()
		//buf := new(bytes.Buffer)
		//buf.ReadFrom(body)

		//accessLogMap["request_post_data"] = c.Request.PostForm.Encode()
		// todo body
		accessLogMap["request_client_ip"] = c.ClientIP()
		ch <- accessLogMap
	}
}
