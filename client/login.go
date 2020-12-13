package client

import (
	"bili/server"
	"bili/tools"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

var (
	loginInfoSleep = time.Tick(time.Second * 3)
	loginInfoStop  = time.Tick(time.Second * 60)
)

func init() {
}

type info struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type navInfo struct {
	info
	Ttl  int64                  `json:"ttl"`
	Data map[string]interface{} `json:"data"`
}

type loginInfo struct {
	info
	States bool        `json:"status"`
	Ts     int64       `json:"ts"`
	Data   interface{} `json:"data"`
}

type loginUrlInfo struct {
	Code   int               `json:"code"`
	States bool              `json:"states"`
	Ts     int64             `json:"ts"`
	Data   map[string]string `json:"data"`
}

func (n navInfo) String() string {
	return fmt.Sprintf("code:%d, message:\"%s\", ttl:%v, data:%v", n.Code, n.Message, n.Ttl, n.Data)
}

func (l loginInfo) String() string {
	return fmt.Sprintf("state:%v, code:%d, message:\"%s\", ts:%d, data:%v", l.States, l.Code, l.Message, l.Ts, l.Data)
}

func (l loginUrlInfo) String() string {
	return fmt.Sprintf("code:%d, state:%v, ts:%d, data:%v", l.Code, l.States, l.Ts, l.Data)
}

func (n *navInfo) ToStruct(body []byte) (err error) {
	err = json.Unmarshal(body, n)
	return
}

func (l *loginInfo) ToStruct(body []byte) (err error) {
	err = json.Unmarshal(body, l)
	return
}

func (l *loginUrlInfo) ToStruct(body []byte) (err error) {
	tools.Log.Debug(string(body))
	err = json.Unmarshal(body, l)
	return
}

func getNavInfo() (ni navInfo, err error) {
	urlStr := "https://api.bilibili.com/x/web-interface/nav"

	body, err := httpClientGet(urlStr)
	if err != nil {
		return
	}
	err = ni.ToStruct(body)
	return
}

func getLoginUrlInfo() (lui loginUrlInfo, err error) {
	urlStr := "https://passport.bilibili.com/qrcode/getLoginUrl"

	body, err := httpClientGet(urlStr)
	if err != nil {
		return
	}
	err = lui.ToStruct(body)
	return
}

func getLoginInfo(oauthKey string) (li loginInfo, err error) {
	urlStr := "https://passport.bilibili.com/qrcode/getLoginInfo"

	data := make(url.Values)
	data["oauthKey"] = []string{oauthKey}
	data["gourl"] = []string{"https://www.bilibili.com/"}

	body, err := httpClientPostForm(urlStr, data)
	if err != nil {
		return
	}
	tools.Log.Debug(string(body))
	err = li.ToStruct(body)
	return
}

func CheckLogin() bool {
	ni, err := getNavInfo()
	if err != nil {
		tools.Log.Fatal(err)
	}
	if ni.Code == 0 {
		tools.Log.Infof("isLogin:%v, name:%s", true, ni.Data["uname"])
		return true
	}
	tools.Log.Warn(ni)
	return false
}

func Login() error {
	if CheckLogin() {
		return nil
	}

	lui, err := getLoginUrlInfo()
	if err != nil {
		tools.Log.Fatal(err)
	}
	if lui.Code != 0 {
		tools.Log.Fatal(lui)
	}
	oauthKey := lui.Data["oauthKey"]
	oauthKeyUrl := lui.Data["url"]

	// 开启二维码验证服务
	server.NewServer(oauthKeyUrl)
	// 弹出二维码
	server.AlertAddress()

	// 循环请求是否验证
forEnd:
	for {
		select {
		case <-loginInfoSleep:
			li, err := getLoginInfo(oauthKey)
			if err != nil {
				tools.Log.Fatal(err)
			}
			if li.States {
				if !CheckLogin() {
					tools.Log.Fatal("login to nav err!")
				}
				Conf.Save()
				return nil
			}
			tools.Log.Infof("login: %s", li)
		case <-loginInfoStop:
			break forEnd
		}
	}

	return fmt.Errorf("请访问:%s, 扫描二维码登录!!!", server.Address)
}
