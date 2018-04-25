package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var friendLogger = logs.GetBeeLogger()
var adapter = "console"
var level = logs.LevelInfo

//Level define
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarn
	LevelNotice
	LevelInfo
	LevelDebug
)

func init() {
	adapter = beego.AppConfig.DefaultString("log::adapter", "console")
	path := beego.AppConfig.DefaultString("log::path", "")
	level = beego.AppConfig.DefaultInt("log::level", logs.LevelInfo)
	async := beego.AppConfig.DefaultInt64("log::async", 0)

	fmt.Printf("init logger: %s, %s, %d, %d\n", adapter, path, level, async)

	logs.SetLogger(adapter, path)
	logs.SetLevel(level)
	logs.EnableFuncCallDepth(false)
	logs.Async(async)

	//initfriendLogger(path, async)
}

func initfriendLogger(path string, async int64) {
	friendLogger = logs.NewLogger()
	friendLogger.SetLogger(adapter, path)
	friendLogger.SetLevel(level)
	// friendLogger.EnableFuncCallDepth(true)
	// friendLogger.SetLogFuncCallDepth(3)
	friendLogger.Async(async)
}

func caller(msg string, depth int) string {
	pc := make([]uintptr, 10)
	runtime.Callers(depth, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	stack := strings.Split(f.Name(), "/")
	name := strings.Split(file, "/")
	return fmt.Sprintf("[%s:%d][%s] %s", name[len(name)-1], line, stack[len(stack)-1], msg)
}

//Debug ...
func Debug(format string, v ...interface{}) {
	if level >= logs.LevelDebug {
		friendLogger.Debug(caller(format, 3), v...)
	}
}

//Info ...
func Info(format string, v ...interface{}) {
	if level >= logs.LevelInfo {
		friendLogger.Info(caller(format, 3), v...)
	}
}

//Notice ...
func Notice(format string, v ...interface{}) {
	if level >= logs.LevelNotice {
		friendLogger.Notice(caller(format, 3), v...)
	}
}

//Warn ...
func Warn(format string, v ...interface{}) {
	if level >= logs.LevelWarn {
		friendLogger.Warn(caller(format, 3), v...)
	}
}

//Error ...
func Error(format string, v ...interface{}) {
	if level >= logs.LevelError {
		friendLogger.Error(caller(format, 3), v...)
	}
}

//Critical ...
func Critical(format string, v ...interface{}) {
	if level >= logs.LevelCritical {
		friendLogger.Critical(caller(format, 3), v...)
	}
}

//LogWithDepth ...
func LogWithDepth(l, d int, f string, v ...interface{}) {
	if level < l {
		return
	}
	switch level {
	case LevelDebug:
		friendLogger.Debug(caller(f, d), v...)

	case LevelInfo:
		friendLogger.Info(caller(f, d), v...)

	case LevelNotice:
		friendLogger.Notice(caller(f, d), v...)

	case LevelWarn:
		friendLogger.Warn(caller(f, d), v...)

	case LevelError:
		friendLogger.Error(caller(f, d), v...)

	case LevelCritical:
	case LevelAlert:
	case LevelEmergency:
		friendLogger.Critical(caller(f, d), v...)

	default:
		friendLogger.Info(caller(f, d), v...)
	}
}

//Adapter ...
func Adapter() string {
	return adapter
}

//Flush ...
func Flush() {
	friendLogger.Flush()
}
