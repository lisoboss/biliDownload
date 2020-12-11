package server

import (
	"context"
	"html/template"
	"net/http"
	"time"
	""
)

func NewServer(key string) {
	stopTime := time.Tick(time.Second * 30)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmpl, err := template.New("qrcode").Parse(Html)
		if err != nil {
			Log.Fatal(err)
		}
		err = tmpl.Execute(writer, key)
		if err != nil {
			Log.Fatal(err)
		}
	})

	server := &http.Server{Addr: ":80"}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			Log.Info(err)
		}
	}()

	go func() {
		<-stopTime
		err := server.Shutdown(context.TODO())
		if err != nil {
			Log.
		}
	}()
}
