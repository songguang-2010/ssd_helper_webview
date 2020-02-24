package file

import (
	"lib/serror"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// OpenFileWrite ...
func OpenFileWrite(logFilePath string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	return logFile, err
}

// OpenFileRead ...
func OpenFileRead(logFilePath string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
	return logFile, err
}

// OpenFileAppend ...
func OpenFileAppend(logFilePath string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return logFile, err
}

// CloseFile ...
func CloseFile(logFile *os.File) {
	err := logFile.Close()
	if err != nil {
		panic(err)
	}
}

// PathExists ...
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CurrentFile ...
func CurrentFile() (string, error) {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return "", serror.New("Can not get current file info")
	}
	return file, nil
}

// CurrentFilePath ..
func CurrentFilePath() (string, error) {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return "", serror.New("Can not get current file info")
	}
	path, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		return "", err
	}
	return path, nil
}

// CurrentExecPath ...
func CurrentExecPath() (string, error) {
	//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	//将\替换成/
	return strings.Replace(dir, "\\", "/", -1), nil
}

// // CurrentExecPath ...
// func CurrentExecPath() (string, error) {
// 	execPath, _ := os.Executable()
// 	path, _ := filepath.EvalSymlinks(execPath)
// 	exDir := filepath.Dir(path)
// 	return exDir, nil
// }
