package main

import (
	"bili/client"
	"bili/tools"
)

func main() {
	defer tools.Log.Close()

	//Log.Fatal(111111111)

	tools.Log.Info("1233333333")
	loginConf := tools.NewConf()
	tools.Log.Infof("Conf: %v", loginConf)
	client.CheckLoginInfo()

	//ab := map[string][]int{"kk":{1,4},"ee":{1,6}}
	//log.Println(ab["kk"])
	//log.Println(ab["ww"])
	//log.Println(ab["ee"])
}
