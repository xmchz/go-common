package gin_util

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

func BindAndValidate(s interface{}, c *gin.Context) (interface{}, bool) {
	b := binding.Default(c.Request.Method, c.ContentType())
	if err := c.ShouldBindWith(s, b); err != nil {
		log.Infof("bind param failed, host: %s, url: %s, err: %s",
			c.Request.RemoteAddr, c.Request.RequestURI, err.Error())
		c.JSON(400, "bind param err")
		return nil, false
	}
	if err := validator.New().Struct(s); err != nil {
		log.Infof("invalid query, host: %s, url: %s, err: %s",
			c.Request.RemoteAddr, c.Request.RequestURI, err.Error())
		c.String(400, "validate param err")
		return nil, false
	}
	return s, true
}
