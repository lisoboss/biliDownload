package tools

import "github.com/robertkrimen/otto"

var (
	jsStr = `
var videoUrlStr, audioUrlStr

function isDel(t) {
    return 1 === parseInt(('000000' + t.toString(2)).substr(-7).split('').reverse().join('').charAt(3));
}

function setUrl(playinfo) {

	playinfo = JSON.parse(playinfo)

	data = playinfo.data
	
	if (data.durl) {
		if (data.durl.length > 0) {
			videoUrlStr = data.durl[0].url
		} else {
			videoUrlStr = ""
		}
		audioUrlStr = ""
	}
	
	if (data.dash) {
		if (data.dash.video.length > 0) {
			videoUrlStr = data.dash.video[0].baseUrl
		} else {
			videoUrlStr = ""
		}
		if (data.dash.audio.length > 0) {
			audioUrlStr = data.dash.audio[0].baseUrl
		} else {
			audioUrlStr = ""
		}
	}
}
`
	jsVm = otto.New()
)

func init() {
	_, err := jsVm.Run(jsStr)
	if err != nil {
		Log.Fatal(err)
	}
}

func IsDel(attr int) bool {
	value, err := jsVm.Call("isDel", nil, attr)
	if err != nil {
		Log.Error(err)
		return true
	}
	boolean, err := value.ToBoolean()
	if err != nil {
		Log.Error(err)
		return true
	}
	return boolean
}

func GetMediaUrl(playInfo string) (videoUrlStr string, audioUrlStr string) {

	_, err := jsVm.Call("setUrl", nil, playInfo)
	if err != nil {
		Log.Error(err)
	}

	videoUrlStrValue, err := jsVm.Get("videoUrlStr")
	if err != nil {
		Log.Error(err)
	}
	videoUrlStr, err = videoUrlStrValue.ToString()
	if err != nil {
		Log.Error(err)
	}
	audioUrlStrValue, err := jsVm.Get("audioUrlStr")
	if err != nil {
		Log.Error(err)
	}
	audioUrlStr, err = audioUrlStrValue.ToString()
	if err != nil {
		Log.Error(err)
	}

	return
}
