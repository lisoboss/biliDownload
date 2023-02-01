package client

import (
	"bili/tools"
	"encoding/json"
	"fmt"
)

type collect struct {
	Id    int    `json:"id"`
	Fid   int    `json:"fid"`
	Mid   int    `json:"mid"`
	Attr  int    `json:"attr"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Upper struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"upper"`
	CoverType  int    `json:"cover_type"`
	Intro      string `json:"intro"`
	Ctime      int    `json:"ctime"`
	Mtime      int    `json:"mtime"`
	State      int    `json:"state"`
	FavState   int    `json:"fav_state"`
	MediaCount int    `json:"media_count"`
	ViewCount  int    `json:"view_count"`
	Type       int    `json:"type"`
	Link       string `json:"link"`
}

type collectInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Count   int       `json:"count"`
		List    []collect `json:"list"`
		HasMore bool      `json:"has_more"`
	} `json:"data"`
}

type collectMedias struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Duration int    `json:"duration"`
	Pubtime  int    `json:"pubtime"`
	Bvid     string `json:"bvid"`
	Upper    struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
	} `json:"upper"`
	CntInfo struct {
		Collect int `json:"collect"`
		Play    int `json:"play"`
	} `json:"cnt_info"`
}

type collectMediasInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Info struct {
			Id         int    `json:"id"`
			SeasonType int    `json:"season_type"`
			Title      string `json:"title"`
			Cover      string `json:"cover"`
			Upper      struct {
				Mid  int    `json:"mid"`
				Name string `json:"name"`
			} `json:"upper"`
			CntInfo struct {
				Collect int `json:"collect"`
				Play    int `json:"play"`
			} `json:"cnt_info"`
			MediaCount int `json:"media_count"`
		} `json:"info"`
		Medias []collectMedias `json:"medias"`
	} `json:"data"`
}

func geCollectListInfo() (collectList []collect, err error) {
	mid := Conf.UpMid
	page := 0
	if mid <= 0 {
		tools.Log.Fatalf("up mid:%f <= 0", mid)
	}
	for true {
		page++
		urlStr := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/folder/collected/list?pn=%d&ps=20&up_mid=%.0f&platform=web&jsonp=jsonp", page, mid)
		body, err := httpClientGet(urlStr)
		if err != nil {
			return collectList, err
		}
		var cInfo collectInfo
		err = json.Unmarshal(body, &cInfo)
		if err != nil {
			return collectList, err
		}
		tools.Log.Debug(cInfo)
		collectList = append(collectList, cInfo.Data.List...)
		if !cInfo.Data.HasMore {
			break
		}
	}
	return
}

func getCollectMediaInfo(c collect) (medias []collectMedias, err error) {
	pn := 1
	count := 0
	for {
		if count >= c.MediaCount {
			break
		}
		urlStr := fmt.Sprintf("https://api.bilibili.com/x/space/fav/season/list?season_id=%d&pn=%d&ps=20&jsonp=jsonp", c.Id, pn)
		body, err := httpClientGet(urlStr)
		if err != nil {
			return medias, err
		}
		var ci collectMediasInfo
		err = json.Unmarshal(body, &ci)
		if err != nil {
			return medias, err
		}
		tools.Log.Debug(ci)
		if len(ci.Data.Medias) <= 0 {
			break
		}
		pn++
		count += len(ci.Data.Medias)
		medias = append(medias, ci.Data.Medias...)
	}
	return
}

func collectStart() {
	defer filter.Save()
	collectList, err := geCollectListInfo()
	if err != nil {
		tools.Log.Fatal(err)
	}

	cn := len(collectList)
	tools.Log.Infof("collectList:%d", cn)
	for _, c := range collectList {
		Conf.AddExcludeCollects(c.Title)
	}
	Conf.Save()
	for i, c := range collectList {
		tools.Log.Info(i, cn, c.Title, c.MediaCount)

		if Conf.ExcludeCollects[c.Title] {
			tools.Log.Info("已排除:", c.Title)
			continue
		}

		medias, err := getCollectMediaInfo(c)
		if err != nil {
			tools.Log.Fatal(err)
		}
		//tools.Log.Debug(medias)
		mn := len(medias)
		tools.Log.Infof("get medias:%d", mn)

		for j, m := range medias {
			tools.Log.Infof("download favList:%d/%d medias:%d/%d => bv:%s name:%s", i+1, cn, j+1, mn, m.Bvid, m.Title)
			if ok := download(c.Title, m.Title, m.Bvid, m.Title, 1); ok {
				tools.Log.Infof("download favList:%d/%d medias:%d/%d successfully", i+1, cn, j+1, mn)
			} else {
				tools.Log.Errorf("download favList:%d/%d medias:%d/%d unsuccessfully", i+1, cn, j+1, mn)
			}
		}
	}
}
