package db

type Filter interface {
	Add(string, string) bool
	Delete(string, string) bool
	Clear(string) bool

	Adds(string, []string) []bool
	Deletes(string, []string) []bool
	ClearAll() bool

	Init()
	Save()
}
