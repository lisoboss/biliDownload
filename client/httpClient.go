package client

import (
	"bili/db"
	"bili/tools"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	Conf         *db.ConfDb
	httpClient   *http.Client
	rangeCompile = regexp.MustCompile(`/(\d+)`)
	sleepChan    = time.Tick(time.Second * 1)
)

func init() {
	Conf = db.NewConf()
	httpClient = &http.Client{}

	httpClient.Jar, _ = cookiejar.New(nil)

	if len(Conf.Cookies) > 0 {
		for key := range Conf.Cookies {
			tools.Log.Debug(key)
			httpClient.Jar.SetCookies(&url.URL{Host: key, Scheme: "https"}, Conf.Cookie(key))
		}
	}
}

func formatBody(resp *http.Response) (body []byte, err error) {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//log.Printf("wrong status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("ERROR status code: %d", resp.StatusCode)
	}
	if rc := resp.Cookies(); len(rc) > 0 {
		Conf.SetCookie(resp.Request.URL, rc)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	return
}

func httpClientPostForm(urlStr string, data url.Values) (body []byte, err error) {
	<-sleepChan
	tools.Log.Debug("post form", urlStr)
	resp, err := httpClient.PostForm(urlStr, data)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return nil, fmt.Errorf("ERROR postform: %v", err)
	}
	return formatBody(resp)
}

func httpClientGet(urlStr string) (body []byte, err error) {
	<-sleepChan
	tools.Log.Debug("get", urlStr)
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return nil, fmt.Errorf("ERROR postform: %v", err)
	}
	return formatBody(resp)
}

func httpClientDownloadByLength(fileOut *os.File, urlStr string, rangeStr string) (lengthMax int, err error) {
	<-sleepChan
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return 0, err
	}
	req.Header["Connection"] = []string{"keep-alive"}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"}
	req.Header["Accept"] = []string{"*/*"}
	req.Header["Origin"] = []string{"https://www.bilibili.com"}
	req.Header["Sec-Fetch-Site"] = []string{"cross-site"}
	req.Header["Sec-Fetch-Mode"] = []string{"cors"}
	req.Header["Sec-Fetch-Dest"] = []string{"empty"}
	req.Header["Referer"] = []string{"https://www.bilibili.com"}
	req.Header["Accept-Encoding"] = []string{"identity"}
	req.Header["Accept-Language"] = []string{"zh-CN,zh;q=0.9,en;q=0.8"}
	req.Header["Range"] = []string{rangeStr}
	resp, err := httpClient.Do(req)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return 0, fmt.Errorf("error download: %v", err)
	}
	defer resp.Body.Close()

	if rangeStr == "bytes=0-10485760" && resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
		//log.Printf("wrong status code: %d", resp.StatusCode)
		contentRangeStr := resp.Header.Get("Content-Range")
		if len(contentRangeStr) <= 0 {
			return 0, err
		}

		rangeStrings := rangeCompile.FindStringSubmatch(contentRangeStr)
		if len(rangeStrings) != 2 {
			return 0, fmt.Errorf("error rangeStrings:%d != 2", len(rangeStrings))
		}

		rexpLengthMax, err := strconv.Atoi(rangeStrings[1])
		if err != nil {
			return 0, err
		}

		return httpClientDownloadByLength(fileOut, urlStr, fmt.Sprintf("bytes=0-%d", rexpLengthMax-1))

	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		//log.Printf("wrong status code: %d", resp.StatusCode)
		text, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("error status code: %d\n%s", resp.StatusCode, text)
	}
	if rc := resp.Cookies(); len(rc) > 0 {
		Conf.SetCookie(resp.Request.URL, rc)
	}

	_, err = io.Copy(fileOut, resp.Body)

	contentRangeStr := resp.Header.Get("Content-Range")
	if len(contentRangeStr) <= 0 {
		return 0, err
	}

	rangeStrings := rangeCompile.FindStringSubmatch(contentRangeStr)
	if len(rangeStrings) != 2 {
		return 0, fmt.Errorf("error rangeStrings:%d != 2", len(rangeStrings))
	}

	return strconv.Atoi(rangeStrings[1])
}

func httpClientDownload(fileOut *os.File, urlStr string) (err error) {
	tools.Log.Debug("download", urlStr)
	length := 1024 * 1024 * 10

	lengthMax := 1
	min := 0
	max := 0
	backBuff := new(bytes.Buffer)
	for min = 0; lengthMax > 0; min = max + 1 {
		max = min + length
		if lengthMax != 1 && max > lengthMax {
			max = lengthMax
		}
		rangeStr := fmt.Sprintf("bytes=%d-%d", min, max)
		fmt.Print(".")
		backBuff.WriteString("\b")
		lengthMax, err = httpClientDownloadByLength(fileOut, urlStr, rangeStr)
		if err != nil {
			tools.Log.Error(lengthMax, rangeStr)
			return err
		}
		lengthMax--
		if lengthMax <= max {
			lengthMax = -1
		}
	}
	fmt.Print(string(backBuff.Bytes()))
	return err
}
