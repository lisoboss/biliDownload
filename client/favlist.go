package client

import (
	"bili/db"
	"bili/tools"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var (
	filter             db.Filter
	fileDir            = "./data"
	filterKey          = "FILEPATH"
	compileSpecialChar = regexp.MustCompile(`[\\/:*?"<>|]`)
	compileMediaInfo   = regexp.MustCompile(`"video":\[({"id".*?),"audio":\[({"id".*?),"support_formats"`)
	compileMediaUrl    = regexp.MustCompile(`"id"[^"]+"baseUrl"[^"]+"([^"]+)"`)
)

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

func getMediaInfo(f fav) (medias []media, err error) {
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

func saveFile(path string, data []byte, createDirBool bool) bool {

	if createDirBool {
		if err := tools.CreateDirFromFilePath(path); err != nil {
			tools.Log.Fatal(err)
		}
	}

	if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
		tools.Log.Error(err)
		return false
	}

	return true
}

func saveFileFromHttp(path string, urlStr string, createDirBool bool) bool {
	if createDirBool {
		if err := tools.CreateDirFromFilePath(path); err != nil {
			tools.Log.Fatal(err)
		}
	}

	fileOut, err := os.Create(path)
	if err != nil {
		tools.Log.Fatal(err)
	}
	defer func() {
		_ = fileOut.Close()
	}()

	err = httpClientDownload(fileOut, urlStr)
	if err != nil {
		tools.Log.Fatal(err)
	}

	return true
}

func download(k string, m media) bool {
	rootPath := filepath.Join(fileDir, k, compileSpecialChar.ReplaceAllString(m.Title, "_"))
	if filter.Exist(filterKey, rootPath) {
		tools.Log.Info("exist:", rootPath)
		return true
	}

	for p := 1; p < m.Page+1; p++ {
		pathStr := filepath.Join(rootPath, fmt.Sprintf("P%d", p))

		if filter.Exist(filterKey, pathStr) {
			tools.Log.Info("exist:", pathStr)
			continue
		}

		urlStr := fmt.Sprintf("https://www.bilibili.com/video/%s?p=%d", m.BvId, p)
		tools.Log.Infof("save > %s", pathStr)
		tools.Log.Infof("download url:%s", urlStr)

		body, err := httpClientGet(urlStr)
		if err != nil {
			tools.Log.Fatal(err)
		}

		result := compileMediaInfo.FindSubmatch(body)
		//tools.Log.Debug(result)
		if len(result) != 3 {
			tools.Log.Warn("filter err regexp result != 1 && result[0] != 3")
			return false
		}

		videoUrlStr := string(compileMediaUrl.FindSubmatch(result[1])[1])
		audioUrlStr := string(compileMediaUrl.FindSubmatch(result[2])[1])

		//tools.Log.Debug(videoUrlStr)
		//tools.Log.Debug(audioUrlStr)

		tools.Log.Info("download video...")
		if ok := saveFileFromHttp(filepath.Join(pathStr, "video.mp4"), videoUrlStr, true); !ok {
			tools.Log.Errorf("save video file %s err", videoUrlStr)
		}
		tools.Log.Info("download audio...")
		if ok := saveFileFromHttp(filepath.Join(pathStr, "audio.mp3"), audioUrlStr, false); !ok {
			tools.Log.Errorf("save audio file %s err", audioUrlStr)
		}
		tools.Log.Info("download info...")
		if ok := saveFile(filepath.Join(pathStr, "info.txt"), []byte(m.Intro), false); !ok {
			tools.Log.Errorf("save info file err")
		}

		if !filter.Add(filterKey, pathStr) {
			tools.Log.Error("filter add err pathStr:", pathStr)
			return false
		}

		filter.Save()
	}

	if !filter.Add(filterKey, rootPath) {
		tools.Log.Error("filter add err rootPath:", rootPath)
		return false
	}

	filter.Save()

	return true
}

func Start() {
	defer filter.Save()
	favList, err := getFavListInfo()
	if err != nil {
		tools.Log.Fatal(err)
	}

	fn := len(favList)
	tools.Log.Infof("favList:%d", fn)
	for i, f := range favList {
		tools.Log.Info(i, fn, f.Title, f.MediaCount)

		if Conf.ExcludeFavorites[f.Title] {
			tools.Log.Info("已排除:", f.Title)
			continue
		}

		medias, err := getMediaInfo(f)
		if err != nil {
			tools.Log.Fatal(err)
		}
		//tools.Log.Debug(medias)
		mn := len(medias)
		tools.Log.Infof("get medias:%d", mn)

		for j, m := range medias {
			tools.Log.Infof("download favList:%d/%d medias:%d/%d => bv:%s name:%s", i+1, fn, j+1, mn, m.BvId, m.Title)
			if ok := download(f.Title, m); ok {
				tools.Log.Infof("download favList:%d/%d medias:%d/%d successfully", i+1, fn, j+1, mn)
			} else {
				tools.Log.Errorf("download favList:%d/%d medias:%d/%d unsuccessfully", i+1, fn, j+1, mn)
			}
		}
	}
}
