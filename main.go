package main

import (
	"github.com/zserge/webview"
	// "github.com/zserge/webview/tree/webview-x"
	// "net/url"
	// "encoding/json"
	"fmt"
	"lib/route"
	// "github.com/spf13/viper"
	"lib/config"
	"lib/file"
	"path/filepath"
	// log "lib/logwrap"
	"controller"
	"lib/log"
	// "lib/serror"
	"middleware"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	// "strconv"
)

// GetPublicPath 获取静态文件服务目录
func GetPublicPath() string {
	execPath, _ := os.Executable()
	path, _ := filepath.EvalSymlinks(execPath)
	exDir := filepath.Dir(path)
	publicDir := http.Dir(exDir + "/app")
	return string(publicDir)
}

// 这里记录所有的应该注册的结构体

// 控制器map
var ControllerMap map[string]interface{}

// 中间件map
var MiddlewareMap map[string]interface{}

//初始化函数
func init() {
	ControllerMap = make(map[string]interface{})
	MiddlewareMap = make(map[string]interface{})

	// 给这两个map赋初始值 每次添加完一条路由或中间件，都要在此处把路由或者中间件注册到这里

	// 注册中间件
	MiddlewareMap["WebMiddleware"] = &middleware.WebMiddleware{}

	// 注册控制器
	ControllerMap["WebController"] = &controller.WebController{}

	//添加路由项
	route.AddRoute(route.RouteItem{
		Method:     "POST",
		Path:       "/login",
		Controller: "WebController",
		Function:   "Login",
	})
}

//服务监听控制器结构
type MyMux struct {
}

//服务监听控制器方法
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fmt.Println(r.Method)
	//返回数据格式是json
	r.ParseForm()
	// fmt.Println("收到客户端请求: ", r.Form)
	// fmt.Println("收到客户端请求: ", r.Form.Get("title"))
	// 这个变量用来标记是否找到了动态路由
	flag := false
	// 每一个http请求都会走到这里，然后在这里，根据请求的URL，为其分配所需要调用的方法
	params := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
	for _, v := range route.RoutesList {
		// 检测路由，根据路由指向需要的数据
		if r.URL.Path == v.Path && r.Method == v.Method {
			// 寻找到了对应路由，无需使用静态服务器
			flag = true

			// 检测该路由中是否存在中间件，如果存在，顺序调用
			for _, m := range v.Middleware {
				// 判断是否注册了这个中间件
				if mid, ok := MiddlewareMap[m]; ok {
					rmid := reflect.ValueOf(mid)
					// 执行中间件，返回values数组
					params = rmid.MethodByName("Handle").Call(params)
					// 判断中间件执行结果，是否还要继续往下走
					str := rmid.Elem().FieldByName("ResString").String()
					if str != "" {
						status := rmid.Elem().FieldByName("Status").Int()
						// 字符串不空，查看状态码，默认返回500错误
						if status == 0 {
							status = 500
						}
						w.WriteHeader(int(status))
						fmt.Fprint(w, str)

						return
					}
				}
			}

			// 检测成功，开始调用方法
			// 获取一个控制器包下的结构体
			// 存在  c为结构体，调用c上挂载的方法
			if d, ok := ControllerMap[v.Controller]; ok {
				reflect.ValueOf(d).MethodByName(v.Function).Call(params)
			}

			// 停止向后执行
			return
		}
	}

	// 如果路由列表中还是没有的话,去静态服务器中寻找
	if !flag {
		// 去静态目录中寻找
		http.ServeFile(w, r, GetPublicPath()+r.URL.Path)
	}

	// http.NotFound(w, r)
	return
}

// func registerController(path string, f func(w http.ResponseWriter, r *http.Request)) {
// 	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
// 		//返回数据格式是json
// 		r.ParseForm()
// 		// fmt.Println("收到客户端请求: ", r.Form)
// 		// fmt.Println("收到客户端请求: ", r.Form.Get("title"))
// 		f(w, r)
// 	})
// }

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

	// http.ListenAndServe(":9090", nil) //设置监听的端口

	ln, err := net.Listen("tcp", "127.0.0.1:39493")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// 设置静态文件目录及服务器
	// fsh := http.FileServer(GetPublicPath())
	// http.Handle("/", http.StripPrefix("/", fsh))

	// registerController("/login", login)
	// registerController("/get-sku-specs", getSkuSpecs)
	// registerController("/get-sku-requests", getSkuRequests)
	// registerController("/get-sku-responses", getSkuResponses)
	// registerController("/get-ssd-orders", getSsdOrders)
	// registerController("/get-aos-orders", getAosOrders)
	// registerController("/get-tps-orders", getTpsOrders)
	// registerController("/get-misc-devices", getMiscDevices)

	mux := &MyMux{}

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
