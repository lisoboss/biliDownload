package client

import (
	"bili/tools"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	fileDir            = "./data"
	filterKey          = "FILEPATH"
	compileSpecialChar = regexp.MustCompile(`[\\/:*?"<>|]`)
	compileMediaInfo   = regexp.MustCompile(`<script>[^<]*?window\.__playinfo__[^{]+({.*?)</\s*script>`)
	compileMediaTitle  = regexp.MustCompile(`<title[^>]+>([^<]+)_哔哩哔哩_bilibili</title`)
)

func saveFile(path string, data []byte, createDirBool bool) bool {
	if createDirBool {
		if err := tools.CreateDirFromFilePath(path); err != nil {
			tools.Log.Fatal(err)
		}
	}

	if err := os.WriteFile(path, data, os.ModePerm); err != nil {
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

func download(k, title, bvId, intro string, page int) bool {
	rootPath := filepath.Join(fileDir, k, compileSpecialChar.ReplaceAllString(title, "_"))
	if filter.Exist(filterKey, rootPath) {
		tools.Log.Info("exist:", rootPath)
		return true
	}

	for p := 1; p < page+1; p++ {
		urlStr := fmt.Sprintf("https://www.bilibili.com/video/%s?p=%d", bvId, p)
		tools.Log.Infof("download url:%s", urlStr)
		body, err := httpClientGet(urlStr)
		if err != nil {
			tools.Log.Fatal(err)
		}

		result := compileMediaTitle.FindSubmatch(body)
		tools.Log.Debug(result)
		if len(result) != 2 {
			tools.Log.Warn("compileMediaTitle err regexp len(result) != 2")
			return false
		}

		pathStr := filepath.Join(rootPath, fmt.Sprintf("P%d_%s", p, strings.Trim(string(result[1]), " ")))
		if filter.Exist(filterKey, pathStr) {
			tools.Log.Info("exist:", pathStr)
			continue
		}
		tools.Log.Infof("save > %s", pathStr)

		result = compileMediaInfo.FindSubmatch(body)
		//tools.Log.Debug(result)
		if len(result) != 2 {
			tools.Log.Warn("compileMediaInfo err regexp len(result) != 2")
			return false
		}

		videoUrlStr, audioUrlStr := tools.GetMediaUrl(string(result[1]))

		//tools.Log.Debug(videoUrlStr)
		//tools.Log.Debug(audioUrlStr)
		var (
			hasVideo  = false
			videoPath = ""
			hasAudio  = false
			audioPath = ""
		)
		if len(videoUrlStr) > 0 {
			tools.Log.Info("download video...")
			videoPath = filepath.Join(pathStr, "video.mp4")
			if ok := saveFileFromHttp(videoPath, videoUrlStr, true); !ok {
				tools.Log.Errorf("save video file %s err", videoUrlStr)
			} else {
				hasVideo = true
			}
		} else {
			tools.Log.Errorf("no video")
		}
		if len(audioUrlStr) > 0 {
			tools.Log.Info("download audio...")
			audioPath = filepath.Join(pathStr, "audio.mp3")
			if ok := saveFileFromHttp(audioPath, audioUrlStr, false); !ok {
				tools.Log.Errorf("save audio file %s err", audioUrlStr)
			} else {
				hasAudio = true
			}
		} else {
			tools.Log.Errorf("no audio")
		}
		if hasVideo && hasAudio {
			tools.Log.Info("merge video audio...")
			if err1 := tools.MediaMerge(videoPath, audioPath, videoPath); err1 != nil {
				tools.Log.Errorf("merge video audio err: %s", err1)
			} else {
				tools.Log.Info("remove audio...")
				if err2 := os.Remove(audioPath); err2 != nil {
					tools.Log.Errorf("remove audio err: %s", err2)
				}
			}
		}

		tools.Log.Info("download info...")
		infoStr := fmt.Sprintf("%s\n%s", urlStr, intro)
		if ok := saveFile(filepath.Join(pathStr, "info.txt"), []byte(infoStr), false); !ok {
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
