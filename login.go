package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

var (
	loginInfoSleep = time.Tick(time.Second * 3)
)

func init() {
}

type loginInfo struct {
	State   bool   `json:"state"`
	Code    int    `json:"code"`
	Data    int    `json:"data"`
	Message string `json:"message"`
	Ts      int16  `json:"ts"`
}

func (l loginInfo) String() string {
	return fmt.Sprintf("state:%v, code:%d, data:%d, message:\"%s\", ts:%d", l.State, l.Code, l.Data, l.Message, l.Ts)
}

func CheckLoginInfo() {
	//Log.Infof("Conf: %v, %v", Conf, *Conf)
	li, err := getLoginInfo()
	if err != nil {
		panic(err)
	}
	Log.Infof("loginInfo: %v", li)

	<- loginInfoSleep

	li, err = getLoginInfo()
	if err != nil {
		panic(err)
	}
	Log.Infof("loginInfo: %v", li)

	//Log.Infof("Conf: %v", *Conf)
	Conf.Save()
}

func getLoginInfo() (li loginInfo, err error) {
	//Log.Infof("httpClient: %v", httpClient)
	urlStr := "https://passport.bilibili.com/qrcode/getLoginInfo"

	data := make(url.Values)
	data["oauthKey"] = []string{"68c3e3b67625f625a719f9c4d47db64c"}
	data["gourl"] = []string{"https://www.bilibili.com/"}

	//Log.Infof("data: %v", data)
	body, err := httpClientPostForm(urlStr, data)
	if err != nil {
		//Log.Infof("wrong request: %v", err)
		return loginInfo{}, err
	}
	//Log.Infof("body: %v", string(body))
	err = json.Unmarshal(body, &li)
	if err != nil {
		panic(err)
	}

	//Log.Infof("loginInfo: %v", loginInfo)

	return
}
