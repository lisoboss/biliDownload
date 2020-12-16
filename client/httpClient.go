package client

import (
	"bili/db"
	"bili/tools"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
