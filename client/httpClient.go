package client

import (
	"bili/tools"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	Conf       *tools.ConfData
	httpClient *http.Client
)

func init() {
	Conf = tools.NewConf()
	httpClient = &http.Client{}

	httpClient.Jar, _ = cookiejar.New(nil)

	if len(Conf.Cookies) > 0 {
		for key, _ := range Conf.Cookies {
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
	resp, err := httpClient.PostForm(urlStr, data)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return nil, fmt.Errorf("ERROR postform: %v", err)
	}
	return formatBody(resp)
}

func httpClientGet(urlStr string) (body []byte, err error) {
	//tools.Log.Debug(urlStr)
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return nil, fmt.Errorf("ERROR postform: %v", err)
	}
	return formatBody(resp)
}
