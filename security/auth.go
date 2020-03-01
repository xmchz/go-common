package security

import (
	"errors"
	"github.com/google/uuid"
	"github.com/xmchz/go-common/cache"
	"net/http"
	"net/url"
)

var (
	Token = "token"
)

type AuthHelper interface {
	GetSubject(*http.Request) (string, error) // get sub (username) from request
	GetObject(*http.Request) (string, error)  // get obj (url) from request
	GetAction(*http.Request) (string, error)  // get act (method) from request
	SaveSubject(string) (string, error)  // save sub, return token
}

func NewDefaultAuthHelper(cache cache.Cache) *DefaultAuthHelper {
	return &DefaultAuthHelper{repo: cache}
}

type DefaultAuthHelper struct {
	repo cache.Cache
}

func (helper *DefaultAuthHelper) GetSubject(req *http.Request) (string, error) {
	cookie, err := req.Cookie(Token)
	if err != nil {
		return "", errors.New("token not found in cookies")
	}
	token, _ := url.QueryUnescape(cookie.Value)
	bytes, err := helper.repo.Get(token)
	if err != nil {
		return "", errors.New("token expired")
	}
	sub := string(bytes)
	_,_ = helper.SaveSubject(sub)
	return sub, nil
}

func (helper *DefaultAuthHelper) GetObject(req *http.Request) (string, error) {
	return req.URL.Path, nil
}

func (helper *DefaultAuthHelper) GetAction(req *http.Request) (string, error) {
	return req.Method, nil
}

func (helper *DefaultAuthHelper) SaveSubject(sub string) (string, error) {
	token, err := helper.repo.FindValue(sub)
	if err !=nil {
		u, _ := uuid.NewRandom()
		token = u.String()
	}
	return token, helper.repo.Set(token, []byte(sub))
}
