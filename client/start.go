package client

import (
	"biliDownload/db"
	"biliDownload/tools"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	filter db.Filter
)

func init() {
	filter = db.NewLocalFilter()
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

func Start() {
	favStart()
	collectStart()
}
func ExitWork() {
	//filter.Save()
}
