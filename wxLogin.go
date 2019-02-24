package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WxLoginReq struct {
	Openid     string `json:"openid"`      //用户唯一标识
	Sessionkey string `json:"session_key"` //会话密钥
	Unionid    string `json:"unionid"`     //用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	Errcode    int    `json:"errcode"`     //错误码
	ErrMsg     string `json:"errMsg"`      //错误信息
}

func getwxLoginResult(code string) (wxInfo WxLoginReq, err error) {
	//https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
	appId := "wxb29a76a9a9b2015b"
	appSecret := "d6cee862f80ae7634adf73db47262a48"
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	resp, err := http.Get(fmt.Sprintf(url, appId, appSecret, code))
	if err != nil {
		return wxInfo, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	err = json.Unmarshal([]byte(body), &wxInfo)
	if err != nil {
		return wxInfo, err
	}
	if wxInfo.Errcode != 0 {
		return wxInfo, errors.New(fmt.Sprintf("code: %d, errmsg: %s", wxInfo.Errcode, wxInfo.ErrMsg))
	}
	return wxInfo, nil
}
