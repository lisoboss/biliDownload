package client

import (
	"bili/db"
	"bili/tools"
	"encoding/json"
	"fmt"
)

var filter db.Filter

func init() {
	filter = db.NewLocalFilter()
}

type fav struct {
	Attr       int     `json:"attr"`
	State      int     `json:"fav_state"`
	Fid        float64 `json:"fid"`
	Id         float64 `json:"id"`
	MediaCount int     `json:"media_count"`
	Mid        float64 `json:"mid"`
	Title      string  `json:"title"`
}

type favListInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Count  int         `json:"count"`
		Season interface{} `json:"season"`
		List   []fav       `json:"list"`
	} `json:"data"`
}

type media struct {
	Attr     int         `json:"attr"`
	BvId     string      `json:"bv_id"`
	BvId2    string      `json:"bvid"`
	CntInfo  interface{} `json:"cnt_info"`
	Cover    string      `json:"cover"`
	Ctime    float64     `json:"ctime"`
	Duration int         `json:"duration"`
	FavTime  float64     `json:"fav_time"`
	Id       float64     `json:"id"`
	Intro    string      `json:"intro"`
	Link     string      `json:"link"`
	Page     int         `json:"page"`
	PubTime  float64     `json:"pubtime"`
	Season   interface{} `json:"season"`
	Title    string      `json:"title"`
	Type     int         `json:"type"`
}

type mediasInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Info   interface{} `json:"info"`
		Medias []media     `json:"medias"`
	} `json:"data"`
}

func (f fav) String() string {
	return fmt.Sprintf("mid:%f, fid:%f, id:%f, state:%d, attr:%d, title:%s, media_count:%d", f.Mid, f.Fid, f.Id, f.State, f.Attr, f.Title, f.MediaCount)
}

func (f favListInfo) String() string {
	return fmt.Sprintf("code:%d, message:%s, ttl:%d, data:%v", f.Code, f.Message, f.Ttl, f.Data)
}

func (m mediasInfo) String() string {
	return fmt.Sprintf("code:%d, message:%s, ttl:%d, mediascunt:%d", m.Code, m.Message, m.Ttl, len(m.Data.Medias))
}

func (f *fav) ToStruct(body []byte) (err error) {
	err = json.Unmarshal(body, f)
	return
}

func (f *favListInfo) ToStruct(body []byte) (err error) {
	err = json.Unmarshal(body, f)
	return
}

func (m *mediasInfo) ToStruct(body []byte) (err error) {
	err = json.Unmarshal(body, m)
	return
}

func getFavListInfo() (favList []fav, err error) {
	mid := Conf.UpMid
	if mid <= 0 {
		tools.Log.Fatalf("up mid:%f <= 0", mid)
	}
	urlStr := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/folder/created/list-all?up_mid=%.0f", mid)
	body, err := httpClientGet(urlStr)
	if err != nil {
		return
	}
	var favListInfo favListInfo
	err = favListInfo.ToStruct(body)
	if err != nil {
		return
	}
	tools.Log.Debug(favListInfo)
	favList = favListInfo.Data.List
	return
}

func getMediaInfo(mediaId float64, mediaCount int) (medias []media, err error) {
	pn := 1
	count := 0
	for {
		if count >= mediaCount {
			break
		}
		urlStr := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/resource/list?media_id=%.0f&pn=%d&ps=20", mediaId, pn)
		body, err := httpClientGet(urlStr)
		if err != nil {
			return medias, err
		}
		var mi mediasInfo
		err = mi.ToStruct(body)
		if err != nil {
			return medias, err
		}
		tools.Log.Debug(mi)
		if len(mi.Data.Medias) <= 0 {
			break
		}
		pn++
		count += len(mi.Data.Medias)
		medias = append(medias, mi.Data.Medias...)
	}
	return medias, err
}

func download() {

}

func Start() {
	defer filter.Save()
	favList, err := getFavListInfo()
	if err != nil {
		tools.Log.Fatal(err)
	}
	n := len(favList)
	tools.Log.Infof("favList:%d", n)

	for i, f := range favList {
		tools.Log.Info(i, n, f.Title, f.MediaCount)
		medias, err := getMediaInfo(f.Id, f.MediaCount)
		if err != nil {
			tools.Log.Fatal(err)
		}
		//tools.Log.Debug(medias)
		tools.Log.Infof("get medias:%d", len(medias))

	}

}
