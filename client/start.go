package client

import (
	"bili/db"
)

var (
	filter db.Filter
)

func init() {
	filter = db.NewLocalFilter()
}

func Start() {
	favListStart()
}
func ExitWork() {
	//filter.Save()
}
