package tools

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Log *_Logger

const (
	logConfigFilePath = "./conf/logs.config.json"
	logConfigStr      = `{
	"logs" : [{
	"name":"default", 
	"data":[
		{"handle":"console", "level":1},
		{"handle":"rotating", "dir":"./log", "filename":"default", "level":0, "maxnum":0, "maxsize":"1MB"}
		]
	}]
}`
)

func initLogConfigFile() error {
	return os.WriteFile(logConfigFilePath, []byte(logConfigStr), os.ModePerm)
}

func init() {
	s, err := os.Stat(logConfigFilePath)

	if err != nil {
		err := CreateDirFromFilePath(logConfigFilePath)
		if err != nil {
			panic(err)
		}
		err = initLogConfigFile()
		if err != nil {
			panic(err)
		}
	} else {
		dir := s.IsDir()
		if dir == true {
			err = initLogConfigFile()
			if err != nil {
				panic(err)
			}
		}
	}

	Log = &_Logger{
		handlers: []Handler{},
	}

	err = NewLogger("default")
	if err != nil {
		log.Fatal(err)
	}
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

const (
	INFINITE int = 0
)

/*
===================

	utils functions

===================
*/
func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// Handler log handlers
type Handler interface {
	SetOutput(w io.Writer)
	Output(calldepth int, s string)
	Outputf(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Flags() int
	SetFlags(flag int)
	SetLevel(level Level)
	Prefix() string
	SetPrefix(prefix string)
	close()
}

type LogHandler struct {
	lg    *log.Logger
	mu    sync.Mutex
	level Level
}

type ConsoleHander struct {
	LogHandler
}

type FileHandler struct {
	LogHandler
	logfile *os.File
}

type RotatingHandler struct {
	LogHandler
	dir      string
	filename string
	maxNum   int
	maxSize  int64
	filetime time.Time
	suffix   int
	logfile  *os.File
}

/*
===================

	json config

===================
*/
type logConfig struct {
	Handle   string `json:"handle"`
	Dir      string `json:"dir"`
	Filename string `json:"filename"`
	Level    int    `json:"level"`
	Maxnum   int    `json:"maxnum"`
	Maxsize  string `json:"maxsize"`
}

type logConfigs struct {
	Name string      `json:"name"`
	Data []logConfig `json:"data"`
}

type configs struct {
	Logs []logConfigs `json:"logs"`
}

var Console, _ = NewConsoleHandler()

func NewConsoleHandler() (*ConsoleHander, error) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	return &ConsoleHander{LogHandler: LogHandler{lg: l}}, nil
}

func NewFileHandler(filepath string) (*FileHandler, error) {
	i := strings.LastIndex(filepath, "\\")
	if -1 == i {
		i = strings.LastIndex(filepath, "/")
		if -1 == i {
			return nil, fmt.Errorf("error filepath: %v", filepath)
		}
	}
	dir := filepath[:i]
	err := os.MkdirAll(dir, 0711)
	if err != nil {
		return nil, err
	}

	logfile, _ := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	l := log.New(logfile, "", log.LstdFlags)
	return &FileHandler{
		LogHandler: LogHandler{lg: l},
		logfile:    logfile,
	}, nil
}

func NewRotatingHandler(dir string, filename string, maxNum int, maxSize int64) (*RotatingHandler, error) {
	if maxNum < 0 {
		return nil, errors.Errorf("maxNum is less than 0")
	}
	err := os.MkdirAll(dir, 0711)
	if err != nil {
		return nil, err
	}
	h := &RotatingHandler{
		dir:      dir,
		filename: filename,
		maxNum:   maxNum,
		maxSize:  maxSize,
		suffix:   0,
	}
	h.newFileData()

	logfile, _ := os.OpenFile(h.generateFileName(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	l := log.New(logfile, "", log.LstdFlags)
	h.LogHandler = LogHandler{lg: l}
	h.logfile = logfile

	if h.isMustRename() {
		h.rename()
	}

	// monitor filesize per second
	go func() {
		timer := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timer.C:
				h.fileCheck()
			}
		}
	}()

	return h, nil
}

func newHandler(lg logConfig) (Handler, error) {
	if "console" == lg.Handle {
		return NewConsoleHandler()
	} else if "file" == lg.Handle {
		return NewFileHandler(lg.Filename)
	} else if "rotating" == lg.Handle {
		l := len(lg.Maxsize)
		if l < 3 {
			return nil, fmt.Errorf("error maxsize:%v", lg.Maxsize)
		}
		unitStr := lg.Maxsize[l-2:]
		maxSizeStr := lg.Maxsize[:l-2]
		maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(unitStr) {
		case "mb":
			maxSize = maxSize * 1024 * 1024
		case "kb":
			maxSize = maxSize * 1024
		default:
			return nil, fmt.Errorf("error maxsize type:%v", unitStr)
		}
		return NewRotatingHandler(lg.Dir, lg.Filename, lg.Maxnum, maxSize)
	}
	return nil, fmt.Errorf("unknown handle:%v", lg.Handle)
}

func NewLogger(name string) error {
	filename := logConfigFilePath
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	c := configs{}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}

	for _, lgs := range c.Logs {
		if name == lgs.Name {
			for _, lg := range lgs.Data {
				if lg.Level < 0 || lg.Level > int(ERROR) {
					Log.Close()
					return fmt.Errorf("unknown log name:%v level:%v", name, lg.Level)
				}
				handler, err := newHandler(lg)
				if err != nil {
					Log.Close()
					return err
				}
				handler.SetLevel(Level(lg.Level))
				handler.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)
				Log.handlers = append(Log.handlers, handler)
			}
		}
	}
	if len(Log.handlers) <= 0 {
		return fmt.Errorf("create logger error:%v", name)
	}

	return nil
}

func (l *LogHandler) Flags() int {
	return l.lg.Flags()
}

func (l *LogHandler) SetFlags(flag int) {
	l.lg.SetFlags(flag)
}

func (l *LogHandler) SetLevel(level Level) {
	l.level = level
}

func (l *LogHandler) Prefix() string {
	return l.lg.Prefix()
}

func (l *LogHandler) SetPrefix(prefix string) {
	l.lg.SetPrefix(prefix)
}

func (l *LogHandler) SetOutput(w io.Writer) {
	l.lg.SetOutput(w)
}

func (l *LogHandler) Output(calldepth int, s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.lg.Output(calldepth, s)
}

func (l *LogHandler) Outputf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.lg.Output(4, fmt.Sprintf(format, v...))
}

func (l *LogHandler) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.Output(4, fmt.Sprintln("debug", v))
	}
}

func (l *LogHandler) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.Outputf("debug ["+format+"]", v...)
	}
}

func (l *LogHandler) Info(v ...interface{}) {
	if l.level <= INFO {
		l.Output(4, fmt.Sprintln("info", v))
	}
}

func (l *LogHandler) Infof(format string, v ...interface{}) {
	if l.level <= INFO {
		l.Outputf("info ["+format+"]", v...)
	}
}

func (l *LogHandler) Warn(v ...interface{}) {
	if l.level <= WARN {
		l.Output(4, fmt.Sprintln("warn", v))
	}
}

func (l *LogHandler) Warnf(format string, v ...interface{}) {
	if l.level <= WARN {
		l.Outputf("warn ["+format+"]", v...)
	}
}

func (l *LogHandler) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.Output(4, fmt.Sprintln("error", v))
	}
}

func (l *LogHandler) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.Outputf("error ["+format+"]", v...)
	}
}

func (l *LogHandler) Panic(v ...interface{}) {
	if l.level <= PANIC {
		l.Output(4, fmt.Sprintln("panic", v))
	}
}

func (l *LogHandler) Panicf(format string, v ...interface{}) {
	if l.level <= PANIC {
		l.Outputf("panic ["+format+"]", v...)
	}
}

func (l *LogHandler) Fatal(v ...interface{}) {
	if l.level <= FATAL {
		l.Output(4, fmt.Sprintln("fatal", v))
	}
}

func (l *LogHandler) Fatalf(format string, v ...interface{}) {
	if l.level <= FATAL {
		l.Outputf("fatal ["+format+"]", v...)
	}
}

func (l *LogHandler) close() {

}

func (h *FileHandler) close() {
	if h.logfile != nil {
		h.mu.Lock()
		defer h.mu.Unlock()
		_ = h.logfile.Close()
	}
}

func (h *RotatingHandler) close() {
	if h.logfile != nil {
		h.mu.Lock()
		defer h.mu.Unlock()
		_ = h.logfile.Close()
	}
}

func (h *RotatingHandler) isMustRename() bool {
	if fileSize(h.generateFileName()) >= h.maxSize {
		return true
	}
	t := time.Now()
	if t.Year() > h.filetime.Year() ||
		t.Year() == h.filetime.Year() && t.Month() > h.filetime.Month() ||
		t.Year() == h.filetime.Year() && t.Month() == h.filetime.Month() && t.Day() > h.filetime.Day() {

		h.newFileData()
		return true
	}
	return false
}

func (h *RotatingHandler) newFileData() {
	h.filetime = time.Now()
	if INFINITE == h.maxNum {
		h.suffix = -1
	} else {
		h.suffix = h.maxNum - 1
	}
}

func (h *RotatingHandler) generateFileName() string {
	filetime := h.filetime.Format("20060102150405")
	if h.suffix <= 0 {
		return fmt.Sprintf("%s/%s.%s.0.log", h.dir, h.filename, filetime)
	}
	return fmt.Sprintf("%s/%s.%s.%d.log", h.dir, h.filename, filetime, h.suffix)
}

func (h *RotatingHandler) rename() {
	if INFINITE == h.maxNum {
		h.suffix = h.suffix + 1
	} else {
		if 0 == (h.suffix+1)%h.maxNum {
			h.suffix = 0
		} else {
			h.suffix = h.suffix%h.maxNum + 1
		}
	}

	newPath := h.generateFileName()
	if isExist(newPath) {
		_ = os.Remove(newPath)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.logfile != nil {
		_ = h.logfile.Close()
	}

	h.logfile, _ = os.OpenFile(newPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	h.lg.SetOutput(h.logfile)
}

func (h *RotatingHandler) fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if h.isMustRename() {
		h.rename()
	}
}

/*
===================

	logger

===================
*/
type _Logger struct {
	handlers []Handler
	mu       sync.Mutex
}

func (l *_Logger) Debug(v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Debug(v...)
	}
}

func (l *_Logger) Debugf(format string, v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Debugf(format, v...)
	}
}

func (l *_Logger) Info(v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Info(v...)
	}
}

func (l *_Logger) Infof(format string, v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Infof(format, v...)
	}
}

func (l *_Logger) Warn(v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Warn(v...)
	}
}

func (l *_Logger) Warnf(format string, v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Warnf(format, v...)
	}
}

func (l *_Logger) Error(v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Error(v...)
	}
}

func (l *_Logger) Errorf(format string, v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Errorf(format, v...)
	}
}

func (l *_Logger) Panic(v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Panic(v...)
	}
	panic(fmt.Sprint(v...))
}

func (l *_Logger) Panicf(format string, v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Panicf(format, v...)
	}
	panic(fmt.Sprintf(format, v...))
}

func (l *_Logger) Fatal(v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Fatal(v...)
	}
	os.Exit(1)
}

func (l *_Logger) Fatalf(format string, v ...interface{}) {
	for i := range l.handlers {
		l.handlers[i].Fatalf(format, v...)
	}
	os.Exit(1)
}

func (l *_Logger) Close() {
	for i := range l.handlers {
		l.handlers[i].close()
	}
	l.handlers = l.handlers[0:0]
}
