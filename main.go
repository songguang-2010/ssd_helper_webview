package main

import (
	"github.com/zserge/webview"
	// "github.com/zserge/webview/tree/webview-x"
	// "net/url"
	// "encoding/json"
	"fmt"
	"lib/route"
	"register"
	// "github.com/spf13/viper"
	"lib/config"
	"lib/file"
	// "path/filepath"
	// log "lib/logwrap"
	// "controller"
	"lib/log"
	// "lib/serror"
	// "middleware"
	"net"
	"net/http"
	"os"
	// "reflect"
	"runtime"
	// "strconv"
)

// GetPublicPath 获取静态文件服务目录
func GetPublicPath() string {
	exDir, _ := file.CurrentExecPath()
	publicDir := http.Dir(exDir + "/app")
	return string(publicDir)
}

//初始化函数
func init() {
	// 加载配置信息
	currentPath, err := file.CurrentFilePath()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = config.Init(currentPath + "/conf")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//初始化中间件
	register.InitMiddleware()
	//初始化控制器
	register.InitController()
	//初始化路由项
	register.InitRouter()

	fmt.Println("current routes: ")
	for _, v := range route.Routes {
		fmt.Println(v.Function)
	}

}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	runtime.GOMAXPROCS(1)

	// log.Info("Task 'sync_order_info' start ...")

	//catch panic error
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Runtime Error: ", err)
		}
	}()

	ln, err := net.Listen("tcp", "127.0.0.1:39493")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// 设置静态文件目录及服务器
	// fsh := http.FileServer(GetPublicPath())
	// http.Handle("/", http.StripPrefix("/", fsh))
	// registerController("/login", login)

	mux := &route.MyMux{}
	mux.SetPublicPath(GetPublicPath())

	go func() {
		// Set up your http server here
		log.Fatal(http.Serve(ln, mux))
	}()

	fmt.Println(ln.Addr().String())

	// webview.Open("Hello", "http://"+ln.Addr().String()+"/index.html", 800, 600, true)
	w := webview.New(webview.Settings{
		URL:       "http://" + ln.Addr().String() + "/",
		Title:     "hello",
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     true,
	})
	w.Run()

}
