package tools

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(i, IsDel(i))
	}
}

func TestName2(t *testing.T) {
	for i := 0; i < 100; i++ {
		url, s := GetMediaUrl(`{
    "code": 0,
    "message": "0",
    "ttl": 1,
    "data": {
        "from": "local",
        "result": "suee",
        "message": "",
        "quality": 80,
        "format": "flv",
        "timelength": 378758,
        "accept_format": "flv,flv720,flv480,flv360",
        "accept_description": ["高清 1080P", "高清 720P", "清晰 480P", "流畅 360P"],
        "accept_quality": [80, 64, 32, 16],
        "video_codecid": 7,
        "seek_param": "start",
        "seek_type": "offset",
        "durl": [{
            "order": 1,
            "length": 378758,
            "size": 22002429,
            "ahead": "E4g=",
            "vhead": "AWQAIP/hABpnZAAgrNlAVgXeaWEAAAMAAQAAAwAeDxgxlgEABWjr7PI8",
            "url": "https://upos-sz-mirrorkodo.bilivideo.com/upgcxcode/69/47/47664769/47664769-1-80.flv?e=ig8euxZM2rNcNbRM7WdVhoM17wUVhwdEto8g5X10ugNcXBlqNxHxNEVE5XREto8KqJZHUa6m5J0SqE85tZvEuENvNo8g2ENvNo8i8o859r1qXg8xNEVE5XREto8GuFGv2U7SuxI72X6fTr859r1qXg8gNEVE5XREto8z5JZC2X2gkX5L5F1eTX1jkXlsTXHeux_f2o859IB_&uipk=5&nbs=1&deadline=1608891300&gen=playurl&os=kodobv&oi=1962279002&trid=6a334fd4c96148a1b0fe0dc6182f153eu&platform=pc&upsig=caa6be1220d3c3b58af869d5386d807e&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,platform&mid=167477322&orderid=0,3&agrr=1&logo=80000000",
            "backup_url": [
                "https://upos-sz-mirrorkodo.bilivideo.com/upgcxcode/69/47/47664769/47664769-1-80.flv?e=ig8euxZM2rNcNbRM7WdVhoM17wUVhwdEto8g5X10ugNcXBlqNxHxNEVE5XREto8KqJZHUa6m5J0SqE85tZvEuENvNo8g2ENvNo8i8o859r1qXg8xNEVE5XREto8GuFGv2U7SuxI72X6fTr859r1qXg8gNEVE5XREto8z5JZC2X2gkX5L5F1eTX1jkXlsTXHeux_f2o859IB_&uipk=5&nbs=1&deadline=1608891300&gen=playurl&os=kodobv&oi=1962279002&trid=6a334fd4c96148a1b0fe0dc6182f153eu&platform=pc&upsig=caa6be1220d3c3b58af869d5386d807e&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,platform&mid=167477322&orderid=1,3&agrr=1&logo=40000000",
                "https://upos-sz-mirrorkodob.bilivideo.com/upgcxcode/69/47/47664769/47664769-1-80.flv?e=ig8euxZM2rNcNbRM7WdVhoM17wUVhwdEto8g5X10ugNcXBlqNxHxNEVE5XREto8KqJZHUa6m5J0SqE85tZvEuENvNo8g2ENvNo8i8o859r1qXg8xNEVE5XREto8GuFGv2U7SuxI72X6fTr859r1qXg8gNEVE5XREto8z5JZC2X2gkX5L5F1eTX1jkXlsTXHeux_f2o859IB_&uipk=5&nbs=1&deadline=1608891300&gen=playurl&os=kodobbv&oi=1962279002&trid=6a334fd4c96148a1b0fe0dc6182f153eu&platform=pc&upsig=f2e5c3c2d516ac43d830c9474b2028d1&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,platform&mid=167477322&orderid=2,3&agrr=1&logo=40000000"
            ]
        }],
        "support_formats": [{
            "quality": 80,
            "format": "flv",
            "new_description": "1080P 高清",
            "display_desc": "1080P",
            "superscript": ""
        }, {
            "quality": 64,
            "format": "flv720",
            "new_description": "720P 高清",
            "display_desc": "720P",
            "superscript": ""
        }, {
            "quality": 32,
            "format": "flv480",
            "new_description": "480P 清晰",
            "display_desc": "480P",
            "superscript": ""
        }, {
            "quality": 16,
            "format": "flv360",
            "new_description": "360P 流畅",
            "display_desc": "360P",
            "superscript": ""
        }]
    },
    "session": "4c917ddbb0053e02d39ea3b0f3e4cc2a"
}`)
		fmt.Println(url, s)
	}
}
