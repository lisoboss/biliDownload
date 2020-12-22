package main

import (
	"bili/client"
	"bili/tools"
)

func main() {
	defer tools.Log.Close()

	tools.Log.Info("程序开始...")

	err := client.Login()
	if err != nil {
		tools.Log.Fatal(err)
	}

	client.Start()
}
