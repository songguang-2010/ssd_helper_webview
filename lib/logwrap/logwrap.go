package logwrap

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"lib/file"
	"lib/log"
	"runtime"
	"strings"
	"time"
)

// type singleton struct{}
// var ins *singleton
// var once sync.Once
// func GetIns() *singleton {
//     once.Do(func(){
//         ins = &singleton{}
//     })
//     return ins
// }

type logContent struct {
	App       string `json:"app"`
	Message   string `json:"message"`
	Context   string `json:"context"`
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Uuid      string `json:"uuid"`
}

// log level
const (
	const_debug    = "DEBUG"
	const_info     = "INFO"
	const_warnning = "WARNNING"
	const_error    = "ERROR"
	const_fatal    = "FATAL"
)

// get log content string by parameter
func initLogContent(msg string, level string, timenow time.Time) string {
	var ok bool
	var file string
	var line int
	_, file, line, ok = runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	datetime := timenow.Format("2006-01-02 15:04:05")

	name := viper.GetString("name")

	c := logContent{
		App:       name,
		Message:   msg,
		Context:   fmt.Sprintf("%s:%d", file, line),
		Level:     level,
		Timestamp: datetime,
		Uuid:      "",
	}

	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// get path of log file
func getLogFile(level string, timenow time.Time) string {
	date := timenow.Format("20060102")
	name := viper.GetString("name")
	return "/var/log/service/" + date + "." + name + "." + strings.ToLower(level) + ".log"
}

// write log content to file
func putLogContent(level string, v ...interface{}) {
	timenow := time.Now()

	c := initLogContent(fmt.Sprint(v...), level, timenow)

	logFilePath := getLogFile(level, timenow)

	logFile, err := file.OpenFileAppend(logFilePath)
	if err != nil {
		panic(err)
	}
	defer file.CloseFile(logFile)

	// n, err1 := io.WriteString(logFile, wireteString) //写入文件(字符串)
	// check(err1)

	logHandler := log.New(logFile, "", log.Llongfile)
	// logHandler.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	logHandler.SetFlags(0)
	logHandler.Println(c)
}

// log wrapper for fatal level
func Fatal(v ...interface{}) {
	putLogContent(const_fatal, v...)
	fmt.Println(fmt.Sprint(v...))
	// os.Exit(1)
}

// log wrapper for error level
func Error(v ...interface{}) {
	putLogContent(const_error, v...)
}

// log wrapper for debug level
func Debug(v ...interface{}) {
	debug := viper.GetString("debug")
	if strings.ToUpper(debug) == "TRUE" {
		putLogContent(const_debug, v...)
	}
}

// log wrapper for warnning level
func Warnning(v ...interface{}) {
	putLogContent(const_warnning, v...)
}

// log wrapper for info level
func Info(v ...interface{}) {
	putLogContent(const_info, v...)
}
