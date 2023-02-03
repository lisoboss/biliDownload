package server

import (
	"biliDownload/tools"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"time"
)

func NewServer(key string, port int) string {
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

	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}

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

	return fmt.Sprintf("http://0.0.0.0:%d", port)
}

func AlertAddress(port int) {
	command := exec.Command("cmd", "/c", "start", fmt.Sprintf("http://127.0.0.1:%d", port))
	_ = command.Run()
}
