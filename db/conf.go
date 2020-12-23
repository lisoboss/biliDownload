package db

import (
	"bili/tools"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	confDbFilePath = "./data/conf.db.json"
	newConfDb      *ConfDb
)

type ConfDb struct {
	Cookies          map[string]map[string]http.Cookie `json:"cookies"`
	UpMid            float64                           `json:"up_mid"`
	ExcludeFavorites map[string]bool                   `json:"exclude_favorites"`
}

func init() {
	initDbFile()
	newConfDb = &ConfDb{}
	newConfDb.Flush()
	if newConfDb.Cookies == nil {
		newConfDb.Cookies = make(map[string]map[string]http.Cookie)
	}
	if len(newConfDb.ExcludeFavorites) <= 0 {
		newConfDb.ExcludeFavorites = map[string]bool{"不下载视频的收藏夹名称1": true, "默认收藏夹": true}
		newConfDb.Save()
	}
	//log.Printf("init newConfDb: %v", newConfDb)
}

func initDbFile() {
	_, err := os.Stat(confDbFilePath)
	if err != nil {
		err = tools.CreateDirFromFilePath(confDbFilePath)
		if err != nil {
			tools.Log.Fatal(err)
		}
	}
}

func NewConf() *ConfDb {
	//log.Printf("newConfDb: %v", newConfDb)
	return newConfDb
}

func (c *ConfDb) Cookie(key string) (cookies []*http.Cookie) {
	for key1, _ := range c.Cookies[key] {
		cookie := c.Cookies[key][key1]
		cookies = append(cookies, &cookie)
	}
	return
}

func (c *ConfDb) SetCookie(url *url.URL, cookies []*http.Cookie) {
	if c.Cookies == nil {
		c.Cookies = make(map[string]map[string]http.Cookie)
	}
	for _, cookie := range cookies {
		key := cookie.Domain
		if key == "" {
			key = url.Host
		}
		if c.Cookies[key] == nil {
			c.Cookies[key] = make(map[string]http.Cookie)
		}
		c.Cookies[key][cookie.Name] = *cookie
	}
}

func (c *ConfDb) Save() {
	// 保存conf
	bytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(confDbFilePath, bytes, 0777)
	if err != nil {
		panic(err)
	}
}

func (c *ConfDb) Read() {
	// 读取conf
	bytes, err := ioutil.ReadFile(confDbFilePath)
	if err != nil {
		c.Save()
		return
	}
	//log.Println(string(bytes))
	_ = json.Unmarshal(bytes, c)
	//log.Println(*c)
}

func (c *ConfDb) Flush() {
	// 刷新conf
	c.Read()
}
