package main

import (
	"bili/client"
	"bili/tools"
)

func main() {
	defer tools.Log.Close()

	//Log.Fatal(111111111)

	tools.Log.Info("程序开始...")
	//conf := tools.NewConf()
	//tools.Log.Debugf("conf: %v", *conf)
	err := client.Login()
	if err != nil {
		tools.Log.Fatal(err)
	}

	//ab := map[string][]int{"kk":{1,4},"ee":{1,6}}
	//log.Println(ab["kk"])
	//log.Println(ab["ww"])
	//log.Println(ab["ee"])
}
