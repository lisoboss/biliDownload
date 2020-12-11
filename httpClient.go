package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	Conf *ConfData
	httpClient *http.Client
)

func init() {
	Conf = NewConf()
	httpClient = &http.Client{}

	httpClient.Jar, _ = cookiejar.New(nil)

	if len(Conf.Cookies) > 0 {
		for key, _ := range Conf.Cookies {
			httpClient.Jar.SetCookies(&url.URL{Host: key, Scheme: "https"}, Conf.Cookie(key))
		}
	}
}

func httpClientPostForm(urlStr string, data url.Values) (body []byte, err error) {
	resp, err := httpClient.PostForm(urlStr, data)
	if err != nil {
		//log.Printf("ERROR: %v", err)
		return nil, fmt.Errorf("ERROR postform: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//log.Printf("wrong status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("ERROR status code: %d", resp.StatusCode)
	}
	if len(resp.Cookies()) > 0 {
		Conf.SetCookie(resp.Request.URL, httpClient.Jar.Cookies(resp.Request.URL))
	}
	body, _ = ioutil.ReadAll(resp.Body)
	return
}