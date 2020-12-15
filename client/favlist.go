package client

import (
	"bili/tools"
	"encoding/json"
	"fmt"
)

type Fav struct {
	Attr       int     `json:"attr"`
	State      int     `json:"fav_state"`
	Fid        float64 `json:"fid"`
	Id         float64 `json:"id"`
	MediaCount int     `json:"media_count"`
	Mid        float64 `json:"mid"`
	Title      string  `json:"title"`
}

type FavListInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Count  int         `json:"count"`
		Season interface{} `json:"season"`
		List   []Fav       `json:"list"`
	} `json:"data"`
}

func (f Fav) String() string {
	return fmt.Sprintf("mid:%f, fid:%f, id:%f, state:%d, attr:%d, title:%s, media_count:%d", f.Mid, f.Fid, f.Id, f.State, f.Attr, f.Title, f.MediaCount)
}

func (f FavListInfo) String() string {
	return fmt.Sprintf("code:%d, message:%s, ttl:%d, data:%v", f.Code, f.Message, f.Ttl, f.Data)
}

func (f *Fav) ToStruct(body []byte) (err error) {
	err = json.Unmarshal(body, f)
	return
}

func (f *FavListInfo) ToStruct(body []byte) (err error) {
	err = json.Unmarshal(body, f)
	return
}

func getFavListInfo() (favList []Fav, err error) {
	mid := Conf.UpMid
	if mid <= 0 {
		tools.Log.Fatalf("up mid:%f <= 0", mid)
	}
	urlStr := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/folder/created/list-all?up_mid=%.0f", mid)
	body, err := httpClientGet(urlStr)
	if err != nil {
		return
	}
	var favListInfo FavListInfo
	err = favListInfo.ToStruct(body)
	if err != nil {
		return
	}
	tools.Log.Debug(favListInfo)
	favList = favListInfo.Data.List
	return
}

func Start() {
	listInfo, err := getFavListInfo()
	if err != nil {
		tools.Log.Fatal(err)
	}
	tools.Log.Debug(listInfo)
}
