package client

import (
	"bili/db"
	"bili/tools"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var (
	Conf       *db.ConfDb
	httpClient *http.Client
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
	tools.Log.Debug("post form", urlStr)
	resp, err := httpClient.PostForm(urlStr, data)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return nil, fmt.Errorf("ERROR postform: %v", err)
	}
	return formatBody(resp)
}

func httpClientGet(urlStr string) (body []byte, err error) {
	tools.Log.Debug("get", urlStr)
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return nil, fmt.Errorf("ERROR postform: %v", err)
	}
	return formatBody(resp)
}

func httpClientDownload(fileOut *os.File, urlStr string) (err error) {
	tools.Log.Debug("download", urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
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
	resp, err := httpClient.Do(req)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return fmt.Errorf("error download: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//log.Printf("wrong status code: %d", resp.StatusCode)
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}
	if rc := resp.Cookies(); len(rc) > 0 {
		Conf.SetCookie(resp.Request.URL, rc)
	}

	_, err = io.Copy(fileOut, resp.Body)

	return err
}
