package db

import (
	"bili/tools"
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	filterDbFilePath = "./data/filter.db.json"
	local            = &LocalFileFilter{}
)

type LocalFileFilter struct {
	Values map[string]map[string]bool
}

func init() {
	local.Init()
}

func NewLocalFilter() *LocalFileFilter {
	return local
}

func (l *LocalFileFilter) Init() {

	bytes, err := ioutil.ReadFile(filterDbFilePath)
	if err != nil {
		local.Save()
		bytes = []byte{}
	}
	err = json.Unmarshal(bytes, l)
	if err != nil {
		tools.Log.Fatal(err)
	}
}

func (l *LocalFileFilter) Save() {
	bytes, err := json.Marshal(l)
	if err != nil {
		tools.Log.Fatal(err)
	}
	err = ioutil.WriteFile(filterDbFilePath, bytes, os.ModePerm)
	if err != nil {
		tools.Log.Fatal(err)
	}
}

func (l *LocalFileFilter) Add(key string, value string) bool {
	if l.Values[key] == nil {
		l.Values[key] = make(map[string]bool)
	}
	if l.Values[key][value] {
		return false
	}

	l.Values[key][value] = true
	return true
}

func (l *LocalFileFilter) Delete(key string, value string) bool {
	if l.Values[key] != nil {
		l.Values[key][value] = false
	}
	return true
}

func (l *LocalFileFilter) Clear(key string) bool {
	if l.Values[key] != nil {
		l.Values[key] = nil
	}
	return true
}

func (l *LocalFileFilter) Adds(key string, values []string) []bool {
	var result []bool
	for _, value := range values {
		result = append(result, l.Add(key, value))
	}
	return result
}

func (l *LocalFileFilter) Deletes(key string, values []string) []bool {
	var result []bool
	for _, value := range values {
		result = append(result, l.Delete(key, value))
	}
	return result
}

func (l *LocalFileFilter) ClearAll() bool {
	l.Values = nil
	return true
}
