package client

import (
	"bili/tools"
	"encoding/json"
	"fmt"
)

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

type favMedia struct {
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

type favMediasInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Info   interface{} `json:"info"`
		Medias []favMedia  `json:"medias"`
	} `json:"data"`
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
	err = json.Unmarshal(body, &favListInfo)
	if err != nil {
		return
	}
	tools.Log.Debug(favListInfo)
	favList = favListInfo.Data.List
	return
}

func getFavMediaInfo(f fav) (medias []favMedia, err error) {
	pn := 1
	count := 0
	for {
		if count >= f.MediaCount {
			break
		}
		urlStr := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/resource/list?media_id=%.0f&pn=%d&ps=20", f.Id, pn)
		body, err := httpClientGet(urlStr)
		if err != nil {
			return medias, err
		}
		var mi favMediasInfo
		err = json.Unmarshal(body, &mi)
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
	return
}

func favStart() {
	defer filter.Save()
	favList, err := getFavListInfo()
	if err != nil {
		tools.Log.Fatal(err)
	}

	fn := len(favList)
	tools.Log.Infof("favList:%d", fn)
	for _, f := range favList {
		Conf.AddExcludeFavorites(f.Title)
	}
	Conf.Save()
	for i, f := range favList {
		tools.Log.Info(i, fn, f.Title, f.MediaCount)

		if Conf.ExcludeFavorites[f.Title] {
			tools.Log.Info("已排除:", f.Title)
			continue
		}

		medias, err := getFavMediaInfo(f)
		if err != nil {
			tools.Log.Fatal(err)
		}
		//tools.Log.Debug(medias)
		mn := len(medias)
		tools.Log.Infof("get medias:%d", mn)

		for j, m := range medias {
			if tools.IsDel(m.Attr) {
				tools.Log.Warnf("warn favList:%d/%d medias:%d/%d => bv:%s name:%s", i+1, fn, j+1, mn, m.BvId, m.Title)
				continue
			}

			tools.Log.Infof("download favList:%d/%d medias:%d/%d => bv:%s name:%s", i+1, fn, j+1, mn, m.BvId, m.Title)
			if ok := download(f.Title, m.Title, m.BvId, m.Intro, m.Page); ok {
				tools.Log.Infof("download favList:%d/%d medias:%d/%d successfully", i+1, fn, j+1, mn)
			} else {
				tools.Log.Errorf("download favList:%d/%d medias:%d/%d unsuccessfully", i+1, fn, j+1, mn)
			}
		}
	}
}
