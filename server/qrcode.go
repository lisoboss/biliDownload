package server

import (
	"bili/tools"
	"context"
	"html/template"
	"net/http"
	"os/exec"
	"time"
)

var Address = "http://localhost"

func NewServer(key string) {
	stopTime := time.Tick(time.Second * 30)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmpl, err := template.New("qrcode").Parse(Html)
		if err != nil {
			tools.Log.Fatal(err)
		}
		err = tmpl.Execute(writer, key)
		if err != nil {
			tools.Log.Fatal(err)
		}
	})

	server := &http.Server{Addr: ":80"}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			tools.Log.Warn(err)
		}
	}()

	go func() {
		<-stopTime
		err := server.Shutdown(context.TODO())
		if err != nil {
			tools.Log.Fatal(err)
		}
	}()
}

func AlertAddress() {
	command := exec.Command("cmd", "/c", "start", Address)
	_ = command.Run()
}
