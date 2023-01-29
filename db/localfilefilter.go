package db

import (
	"bili/tools"
	"encoding/json"
	"os"
)

var (
	filterDbFilePath = "./conf/filter.db.json"
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
	_, err := os.Stat(filterDbFilePath)
	if err != nil {
		local.Save()
	}

	bytes, err := os.ReadFile(filterDbFilePath)
	if err != nil {
		tools.Log.Fatal(err)
	}
	err = json.Unmarshal(bytes, l)
	if err != nil {
		tools.Log.Fatal(err)
	}

	if l.Values == nil {
		l.Values = make(map[string]map[string]bool)
	}
}

func (l *LocalFileFilter) Save() {
	bytes, err := json.Marshal(l)
	if err != nil {
		tools.Log.Fatal(err)
	}
	err = os.WriteFile(filterDbFilePath, bytes, os.ModePerm)
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
		l.Values[key] = make(map[string]bool)
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
	l.Values = make(map[string]map[string]bool)
	return true
}

func (l *LocalFileFilter) Exist(key string, value string) bool {
	if l.Values[key] == nil {
		return false
	}
	return l.Values[key][value]
}
