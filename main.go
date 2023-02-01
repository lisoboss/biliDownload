package main

import (
	"bili/client"
	"bili/tools"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer tools.Log.Close()

	listenSignal()
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

func listenSignal() {
	//创建监听退出chan
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit...", s)
				client.ExitWork()
				os.Exit(0)
			case syscall.SIGUSR1:
				fmt.Println("usr1 signal", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2 signal", s)
			default:
				fmt.Println("other signal", s)
			}
		}
	}()
}
