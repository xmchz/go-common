package gin_middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/xmchz/go-common/security"
	"github.com/xmchz/go-common/util"
)

var CurrentSubjectKey = "current_subject"

func Authenticate(helper security.AuthHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		if subject, err := helper.GetSubject(c.Request); err != nil {
			c.AbortWithStatusJSON(401, util.BaseResp{Message: "need authentication"})
		} else {
			c.Set(CurrentSubjectKey, subject)
			c.Next()
		}
	}
}

func Authorize(e *casbin.Enforcer, helper security.AuthHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, _ := helper.GetSubject(c.Request)
		obj, _ := helper.GetObject(c.Request)
		act, _ := helper.GetAction(c.Request)
		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(500, util.BaseResp{Message: "error occurred when authorizing"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, util.BaseResp{Message: "forbidden"})
			return
		}
		c.Next()
	}
}

