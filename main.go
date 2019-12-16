package main

import (
	"github.com/zserge/webview"
	// "github.com/zserge/webview/tree/webview-x"
	// "net/url"
	"encoding/json"
	"fmt"
	// "github.com/spf13/viper"
	"lib/config"
	"lib/file"
	// log "lib/logwrap"
	"lib/log"
	"lib/serror"
	ssd_aos "model/aos"
	ssd_order "model/order"
	"model/sku"
	ssd_tps "model/tps"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func get_tps_orders(w http.ResponseWriter, r *http.Request) {
	orderNo := string(r.Form.Get("order_no"))

	//实例化数据模型
	orderInfoModel, err := ssd_tps.CreateOrderInfo()
	serror.Check(err)
	defer func() {
		errClean := orderInfoModel.CloseDB()
		serror.Check(errClean)
	}()

	err = orderInfoModel.SetOrderNo(orderNo)
	serror.Check(err)

	// log.Info("date:", dateCurrent)

	rows, err := orderInfoModel.GetList(orderNo)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_tps.OrderInfoRecord, 0)

	//循环处理
	pubCnt := 0
	for rows.Next() {
		specInfoStruct, err := orderInfoModel.ScanRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, specInfoStruct)
		pubCnt += 1
	}

	response(w, pubArr)
}

func get_aos_orders(w http.ResponseWriter, r *http.Request) {
	phone := string(r.Form.Get("phone"))
	dateCurrent := string(r.Form.Get("date"))
	shopName := ""

	//实例化数据模型
	orderInfoModel, err := ssd_aos.CreateOrderInfo()
	serror.Check(err)
	defer func() {
		errClean := orderInfoModel.CloseDB()
		serror.Check(errClean)
	}()

	err = orderInfoModel.SetCurrentDate(dateCurrent)
	serror.Check(err)

	// log.Info("date:", dateCurrent)

	rows, err := orderInfoModel.GetList(shopName, phone)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_aos.OrderInfoRecord, 0)

	//循环处理
	pubCnt := 0
	for rows.Next() {
		specInfoStruct, err := orderInfoModel.ScanRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, specInfoStruct)
		pubCnt += 1
	}

	response(w, pubArr)
}

func get_ssd_orders(w http.ResponseWriter, r *http.Request) {
	phone := string(r.Form.Get("phone"))
	dateCurrent := string(r.Form.Get("date"))
	shopName := ""

	//实例化数据模型
	orderInfoModel, err := ssd_order.CreateOrderInfo()
	serror.Check(err)
	defer func() {
		errClean := orderInfoModel.CloseDB()
		serror.Check(errClean)
	}()

	err = orderInfoModel.SetCurrentDate(dateCurrent)
	serror.Check(err)

	// log.Info("date:", dateCurrent)

	rows, err := orderInfoModel.GetList(shopName, phone)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_order.OrderInfoRecord, 0)

	//循环处理
	pubCnt := 0
	for rows.Next() {
		specInfoStruct, err := orderInfoModel.ScanRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, specInfoStruct)
		pubCnt += 1
	}

	response(w, pubArr)
}

func get_spec_list(w http.ResponseWriter, r *http.Request) {
	shop_name := string(r.Form.Get("shop_name"))

	//要处理的数据的日期，格式：2019-10-12
	dateCurrent := "2019-12-02"

	//实例化数据模型
	specInfoModel, err := sku.CreateSpecInfo()
	serror.Check(err)
	defer func() {
		errClean := specInfoModel.CloseDB()
		serror.Check(errClean)
	}()

	err = specInfoModel.SetCurrentDate(dateCurrent)
	serror.Check(err)

	// log.Info("date:", dateCurrent)

	rows, err := specInfoModel.GetListByShopName(shop_name, 500)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]sku.SpecInfoRecord, 0)

	//循环处理
	pubCnt := 0
	for rows.Next() {
		specInfoStruct, err := specInfoModel.ScanRow(rows)
		serror.Check(err)

		//计入当前处理数组
		// pubArr[pubCnt] = specInfoStruct
		pubArr = append(pubArr, specInfoStruct)
		pubCnt += 1
	}

	// if pubCnt != 0 {
	// 	pubArrSub := make([]sku.SpecInfoRecord, pubCnt)
	// 	copy(pubArrSub, pubArr)
	// 	//重新初始化数组切片
	// 	pubArr = make([]sku.SpecInfoRecord, 50)
	// }

	// return pubArr, nil
	response(w, pubArr)
}

func registerController(path string, f func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		//返回数据格式是json
		r.ParseForm()
		// fmt.Println("收到客户端请求: ", r.Form)
		// fmt.Println("收到客户端请求: ", r.Form.Get("title"))
		f(w, r)
	})
}

func response(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := string(b)

	// w.Header().Add("Last-Modified", "Thu, 18 Jun 2015 10:24:27 GMT")
	// w.Header().Add("Accept-Ranges", "bytes")
	// w.Header().Add("E-Tag", "55829c5b-17")
	// w.Header().Add("Server", "golang-http-server")
	// w.Write([]byte("<h1>\nHello world!\n</h1>\n"))
	w.Header().Set("Connection", "keep-alive")
	// w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	fmt.Fprint(w, body)
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
	// 加载配置信息
	current_path, err := file.CurrentFilePath()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = config.Init(current_path + "/conf")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:39493")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	// 设置静态目录
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	fsh := http.FileServer(http.Dir(exPath + "/app"))
	http.Handle("/", http.StripPrefix("/", fsh))

	registerController("/get-specs", get_spec_list)
	registerController("/get-ssd-orders", get_ssd_orders)
	registerController("/get-aos-orders", get_aos_orders)
	registerController("/get-tps-orders", get_tps_orders)

	go func() {
		// Set up your http server here
		log.Fatal(http.Serve(ln, nil))
	}()
	fmt.Println(ln.Addr().String())
	webview.Open("Hello", "http://"+ln.Addr().String()+"/index.html", 800, 600, true)

}

// func index(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello golang http!")
// }
