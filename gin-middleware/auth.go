package gin_middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/xmchz/go-common/gin-util"
	"github.com/xmchz/go-common/util"
)

func Authenticate(helper gin_util.AuthHelper) gin.HandlerFunc {

	return func(c *gin.Context) {
		if sub, err := helper.GetSubject(c); err != nil {
			c.AbortWithStatusJSON(401, util.BaseResp{Message: "user hasn't logged in yet"})
		} else {
			_ = helper.SaveSubject(c, sub) // save sub in context, and refresh sub in cache
			c.Next()
		}
	}
}

func Authorize(e *casbin.Enforcer, helper gin_util.AuthHelper) gin.HandlerFunc {

	return func(c *gin.Context) {
		sub, _ := helper.GetSubject(c)
		obj, _ := helper.GetObject(c)
		act, _ := helper.GetAction(c)
		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(500, util.BaseResp{Message: "error occurred when authorizing user"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, util.BaseResp{Message: "forbidden"})
			return
		}
		c.Next()
	}
}

