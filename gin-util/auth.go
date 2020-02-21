package gin_util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xmchz/go-common/util"
)

type AuthHelper interface {
	GetSubject(*gin.Context) (string, error)
	GetObject(*gin.Context) (string, error)
	GetAction(*gin.Context) (string, error)
	SaveSubject(*gin.Context, string) error
}

func NewDefaultAuthHelper(cache util.Cache) *DefaultAuthHelper {
	return &DefaultAuthHelper{cache: cache}
}

type DefaultAuthHelper struct {
	cache util.Cache // token to username(current_subject)
}

func (h *DefaultAuthHelper) GetSubject(c *gin.Context) (string, error) {
	// from ctx
	if sub, ok := c.Get("current_subject"); ok {
		return sub.(string), nil
	}

	// from cookie
	token, _ := c.Cookie("token")
	if bytes, err := h.cache.Get(token) ; err == nil {
		return string(bytes), nil
	}

	return "", errors.New("fail to get subject from ctx or cookies")
}

func (h *DefaultAuthHelper) GetObject(c *gin.Context) (string, error) {
	return c.Request.URL.Path, nil
}

func (h *DefaultAuthHelper) GetAction(c *gin.Context) (string, error) {
	return c.Request.Method, nil
}

func (h *DefaultAuthHelper) SaveSubject(c *gin.Context, sub string) error {
	c.Set("current_subject", sub)	// save in ctx, for c.next
	token, _ := h.cache.FindValue(sub)
	return h.cache.Set(token, []byte(sub)) // save in cache
}