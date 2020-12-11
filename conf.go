package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	rootPath = "./conf_data.json"
	newConfData *ConfData
)

type P struct {
	Title        string `json:"title"`
	SavePath     string `json:"save_path"`
	DownloadPath string `json:"download_path"`
}

type F struct {
	Title string `json:"title"`
	Ps    []P
}

type J struct {
	Title string `json:"title"`
	Fs    []F
}

type ConfData struct {
	Cookies map[string][]http.Cookie `json:"cookies"`
	Js      []J
}

func init() {
	newConfData = &ConfData{}
	newConfData.Flush()
	if newConfData.Cookies == nil {
		newConfData.Cookies = make(map[string][]http.Cookie)
	}
	//log.Printf("init newConfData: %v", newConfData)
}

func NewConf() *ConfData {
	//log.Printf("newConfData: %v", newConfData)
	return newConfData
}

func (c *ConfData) Cookie(key string) (cookies []*http.Cookie) {
	for _, cookie := range c.Cookies[key] {
		cookies = append(cookies, &cookie)
	}
	return
}

func (c *ConfData) SetCookie(u *url.URL, cookies []*http.Cookie) {
	key := u.Host
	if c.Cookies == nil {
		c.Cookies = make(map[string][]http.Cookie)
	}
	for _, cookie := range cookies {
		c.Cookies[key] = append(c.Cookies[key], *cookie)
	}
}

func (c *ConfData) Save() {
	// 保存conf
	bytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(rootPath, bytes, 0777)
	if err != nil {
		panic(err)
	}
}

func (c *ConfData) Read() {
	// 读取conf
	bytes, err := ioutil.ReadFile(rootPath)
	if err != nil {
		c.Save()
		return
	}
	//log.Println(string(bytes))
	_ = json.Unmarshal(bytes, c)
	//log.Println(*c)
}

func (c *ConfData) Flush() {
	// 刷新conf
	c.Read()
}
