package client

import (
	"biliDownload/tools"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	urlStr := "https://223-166-92-250.mcdn.bilivideo.cn:480/upgcxcode/18/21/167382118/167382118-1-30080.m4s?expires=1608789802&platform=pc&ssig=xDetDg-z7g127O7vgnv33A&oi=989301972&trid=1d2ff7d412b44cf9acc4e45c9522f98fu&nfc=1&nfb=maPYqpoel5MI3qOUX6YpRA==&mcdnid=8000016&mid=167477322&orderid=0,3&agrr=0&logo=A0000080"
	fileOut, err := os.Create("./data.out")
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
}
