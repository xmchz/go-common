package oauth

import (
	"fmt"
	"github.com/xmchz/go-common/util"
	"net/http"
	"net/url"
	"strings"
)

type GithubTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubUserInfoResp struct {
	Username  string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
}

func RequestGithubToken(tokenUrl, clientId, clientSecret, authCode string) (string, error){
	postValue := url.Values{
		"client_id":     []string{clientId},
		"client_secret": []string{clientSecret},
		"code":          []string{authCode},
	}
	postStr := postValue.Encode()
	tokenRequest, err := http.NewRequest("POST", tokenUrl, strings.NewReader(postStr))
	if err != nil {
		return "", fmt.Errorf("new tokenRequest err: %s", err.Error())
	}
	tokenRequest.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(tokenRequest)
	if err != nil {
		return "", fmt.Errorf("do tokenRequest err: %s", err.Error())

	}
	gTokenResp := new(GithubTokenResp)
	if err := util.GetBodyAsStruct(resp, gTokenResp); err != nil {
		return "", fmt.Errorf("read token resp body err: %s", err.Error())
	}
	return gTokenResp.AccessToken, nil
}


func RequestGithubUserInfo(userInfoUrl, token, tokenType string) (GithubUserInfoResp, error){
	userInfoRequest, err := http.NewRequest("GET", userInfoUrl, nil)
	if err != nil {
		return GithubUserInfoResp{}, fmt.Errorf("new userInfoRequest err: %s", err.Error())

	}
	userInfoRequest.Header.Add("Authorization",	fmt.Sprintf("%s %s", tokenType, token))
	resp, err := http.DefaultClient.Do(userInfoRequest)

	gUserInfoResp := new(GithubUserInfoResp)
	if err = util.GetBodyAsStruct(resp, gUserInfoResp); err!= nil {
		return GithubUserInfoResp{}, fmt.Errorf("read user info resp body err: %s", err.Error())
	}
	return *gUserInfoResp, nil
}