package main

import (
	"bili/client"
	"bili/tools"
	"time"
)

func main() {
	defer tools.Log.Close()

	tools.Log.Info("程序开始...")

	err := client.Login()
	if err != nil {
		tools.Log.Fatal(err)
	}

	startSleep := time.Tick(time.Minute * 3)

	for {
		client.Start()
		<-startSleep
	}

}
