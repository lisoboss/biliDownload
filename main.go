package main

func main() {
	defer Log.Close()

	//Log.Fatal(111111111)

	Log.Info("1233333333")
	loginConf := NewConf()
	Log.Infof("Conf: %v", loginConf)
	CheckLoginInfo()

	//ab := map[string][]int{"kk":{1,4},"ee":{1,6}}
	//log.Println(ab["kk"])
	//log.Println(ab["ww"])
	//log.Println(ab["ee"])
}
