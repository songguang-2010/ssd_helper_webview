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
	ssd_misc "model/misc"
	ssd_order "model/order"
	ssd_sku "model/sku"
	ssd_stat "model/stat"
	ssd_tps "model/tps"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	// "strconv"
)

func getSkuResponses(w http.ResponseWriter, r *http.Request) {
	shopNo := string(r.Form.Get("shop_no"))
	shopName := string(r.Form.Get("shop_name"))
	searchNoArr := make([]string, 0)
	dateCurrent := string(r.Form.Get("date_response"))

	//实例化数据模型
	shopModel, err := ssd_stat.CreateJwdShopCache()
	serror.Check(err)
	defer func() {
		errClean := shopModel.CloseDB()
		serror.Check(errClean)
	}()

	if shopName != "" {
		//根据门店名称查询门店编码
		rowsShop, err := shopModel.GetListByName(shopName)
		serror.Check(err)
		defer rowsShop.Close()

		for rowsShop.Next() {
			shopRecord, err := shopModel.ScanRow(rowsShop)
			serror.Check(err)
			//计入当前处理数组
			searchNoArr = append(searchNoArr, shopRecord.ShopNo)
		}
	}

	if shopNo != "" {
		searchNoArr = append(searchNoArr, shopNo)
	}

	//实例化数据模型
	goodsPurchaseModel, err := ssd_sku.CreatePurchaseInfo()
	serror.Check(err)
	defer func() {
		errClean := goodsPurchaseModel.CloseDB()
		serror.Check(errClean)
	}()

	rows, err := goodsPurchaseModel.GetResponseList(searchNoArr, dateCurrent)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_sku.PurchaseInfoRecord, 0)

	//循环处理
	for rows.Next() {
		oRecord, err := goodsPurchaseModel.ScanRow(rows)
		serror.Check(err)
		oRecord.Store_code = goodsPurchaseModel.ShopNoUnfilled(oRecord.Store_code)
		//计入当前处理数组
		pubArr = append(pubArr, oRecord)
	}

	//获取门店名称
	if len(pubArr) != 0 {
		shopNoArr := make([]string, 0)
		for _, value := range pubArr {
			shopNoArr = append(shopNoArr, value.Store_code)
		}

		rowsShopList, err := shopModel.GetListByShopNoArr(shopNoArr)
		serror.Check(err)
		defer rowsShopList.Close()

		//初始化shop map
		shopMap := make(map[string]string)

		for rowsShopList.Next() {
			oRecord, err := shopModel.ScanRow(rowsShopList)
			serror.Check(err)
			//计入当前处理数组
			shopMap[oRecord.ShopNo] = oRecord.ShopName
		}

		fmt.Println(shopNoArr)

		fmt.Println(shopMap)

		for key, value := range pubArr {
			pubArr[key].Store_name = shopMap[value.Store_code]
		}
	}

	response(w, pubArr)
}

func getSkuRequests(w http.ResponseWriter, r *http.Request) {
	shopNo := string(r.Form.Get("shop_no"))
	shopName := string(r.Form.Get("shop_name"))
	searchNoArr := make([]string, 0)
	dateCurrent := string(r.Form.Get("date_response"))

	//实例化数据模型
	shopModel, err := ssd_stat.CreateJwdShopCache()
	serror.Check(err)
	defer func() {
		errClean := shopModel.CloseDB()
		serror.Check(errClean)
	}()

	if shopName != "" {
		//根据门店名称查询门店编码
		rowsShop, err := shopModel.GetListByName(shopName)
		serror.Check(err)
		defer rowsShop.Close()

		for rowsShop.Next() {
			shopRecord, err := shopModel.ScanRow(rowsShop)
			serror.Check(err)
			//计入当前处理数组
			searchNoArr = append(searchNoArr, shopRecord.ShopNo)
		}
	}

	if shopNo != "" {
		searchNoArr = append(searchNoArr, shopNo)
	}

	//实例化数据模型
	goodsPurchaseModel, err := ssd_sku.CreatePurchaseInfo()
	serror.Check(err)
	defer func() {
		errClean := goodsPurchaseModel.CloseDB()
		serror.Check(errClean)
	}()

	rows, err := goodsPurchaseModel.GetRequestList(searchNoArr, dateCurrent)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_sku.PurchaseInfoRecord, 0)

	//循环处理
	for rows.Next() {
		oRecord, err := goodsPurchaseModel.ScanRow(rows)
		serror.Check(err)
		oRecord.Store_code = goodsPurchaseModel.ShopNoUnfilled(oRecord.Store_code)
		//计入当前处理数组
		pubArr = append(pubArr, oRecord)
	}

	// type resRecord struct {
	// 	ID           int    `json:"id"`
	// 	StoreCode    string `json:"store_code"`
	// 	StoreName    string `json:"store_name"`
	// 	ProdCode     string `json:"prod_code"`
	// 	ProdName     string `json:"prod_name"`
	// 	ProdNumber   string `json:"prod_number"`
	// 	DateRequest  string `json:"date_request"`
	// 	DateResponse string `json:"date_response"`
	// }

	// //初始化返回数组
	// resArr := make([]resRecord, 0)

	//获取门店名称
	if len(pubArr) != 0 {
		shopNoArr := make([]string, 0)
		for _, value := range pubArr {
			shopNoArr = append(shopNoArr, value.Store_code)
		}

		rowsShopList, err := shopModel.GetListByShopNoArr(shopNoArr)
		serror.Check(err)
		defer rowsShopList.Close()

		//初始化shop map
		shopMap := make(map[string]string)

		for rowsShopList.Next() {
			oRecord, err := shopModel.ScanRow(rowsShopList)
			serror.Check(err)
			//计入当前处理数组
			shopMap[oRecord.ShopNo] = oRecord.ShopName
		}

		fmt.Println(shopNoArr)

		fmt.Println(shopMap)

		for key, value := range pubArr {
			// shopNo := goodsPurchaseModel.ShopNoUnfilled(value.Store_code)
			pubArr[key].Store_name = shopMap[value.Store_code]
			// resArr = append(resArr, resRecord{
			// 	ID:           value.ID,
			// 	StoreCode:    value.StoreCode,
			// 	StoreName:    shopMap[value.StoreCode],
			// 	ProdCode:     value.ProdCode,
			// 	ProdName:     value.ProdName,
			// 	ProdNumber:   value.ProdNumber,
			// 	DateRequest:  value.DateRequest,
			// 	DateResponse: DateResponse,
			// })
		}
	}

	response(w, pubArr)
}

func getMiscDevices(w http.ResponseWriter, r *http.Request) {
	shopNo := string(r.Form.Get("shop_no"))
	shopName := string(r.Form.Get("shop_name"))
	searchNoArr := make([]string, 0)
	appVersion := string(r.Form.Get("app_version"))
	// canary := string(r.Form.Get("canary"))

	// isCanary, err := strconv.Atoi(canary)
	// serror.Check(err)

	//实例化数据模型
	shopModel, err := ssd_stat.CreateJwdShopCache()
	serror.Check(err)
	defer func() {
		errClean := shopModel.CloseDB()
		serror.Check(errClean)
	}()

	if shopName != "" {
		//根据门店名称查询门店编码
		rowsShop, err := shopModel.GetListByName(shopName)
		serror.Check(err)
		defer rowsShop.Close()

		for rowsShop.Next() {
			shopRecord, err := shopModel.ScanRow(rowsShop)
			serror.Check(err)
			//计入当前处理数组
			searchNoArr = append(searchNoArr, shopRecord.ShopNo)
		}
	}

	if shopNo != "" {
		searchNoArr = append(searchNoArr, shopNo)
	}

	//灰度设备id列表
	deviceIDMap := make(map[int]int)

	//实例化数据模型
	canaryDeviceModel, err := ssd_misc.CreateCanaryDevice()
	serror.Check(err)
	defer func() {
		errClean := canaryDeviceModel.CloseDB()
		serror.Check(errClean)
	}()
	//查询灰度设备id列表
	rowsDeviceCanary, err := canaryDeviceModel.GetList()
	serror.Check(err)
	defer rowsDeviceCanary.Close()

	for rowsDeviceCanary.Next() {
		canaryRecord, err := canaryDeviceModel.ScanRow(rowsDeviceCanary)
		serror.Check(err)
		//计入当前灰度设备列表
		deviceIDMap[canaryRecord.DeviceID] = canaryRecord.DeviceID
	}

	//实例化数据模型
	deviceModel, err := ssd_misc.CreateDevice()
	serror.Check(err)
	defer func() {
		errClean := deviceModel.CloseDB()
		serror.Check(errClean)
	}()

	rows, err := deviceModel.GetList(searchNoArr, appVersion)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_misc.Record, 0)

	//循环处理
	for rows.Next() {
		oRecord, err := deviceModel.ScanRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, oRecord)
	}

	type resRecord struct {
		ID          int    `json:"id"`
		ShopNo      string `json:"shop_no"`
		ShopName    string `json:"shop_name"`
		AppVersion  string `json:"app_version"`
		SerialNo    string `json:"serial_no"`
		NetworkType string `json:"network_type"`
		AppEnv      string `json:"app_env"`
		IsCanary    string `json:"is_canary"`
		UpdateTime  string `json:"update_time"`
		CreateTime  string `json:"create_time"`
	}

	//初始化返回数组
	resArr := make([]resRecord, 0)

	//获取门店名称
	if len(pubArr) != 0 {
		shopNoArr := make([]string, 0)
		for _, value := range pubArr {
			shopNoArr = append(shopNoArr, value.ShopNo)
		}

		rowsShopList, err := shopModel.GetListByShopNoArr(shopNoArr)
		serror.Check(err)
		defer rowsShopList.Close()

		//初始化shop map
		shopMap := make(map[string]string)

		for rowsShopList.Next() {
			oRecord, err := shopModel.ScanRow(rowsShopList)
			serror.Check(err)
			//计入当前处理数组
			shopMap[oRecord.ShopNo] = oRecord.ShopName
		}

		for _, value := range pubArr {
			isCanary := "0"
			if deviceIDMap[value.ID] != 0 {
				isCanary = "1"
			}

			resArr = append(resArr, resRecord{
				ID:          value.ID,
				ShopNo:      value.ShopNo,
				ShopName:    shopMap[value.ShopNo],
				AppVersion:  value.AppVersion,
				SerialNo:    value.SerialNo,
				NetworkType: value.NetworkType,
				AppEnv:      value.AppEnv,
				IsCanary:    isCanary,
				UpdateTime:  value.UpdateTime,
				CreateTime:  value.CreateTime,
			})
		}
	}

	response(w, resArr)
}

func getTpsOrders(w http.ResponseWriter, r *http.Request) {
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
	// pubCnt := 0
	for rows.Next() {
		specInfoStruct, err := orderInfoModel.ScanRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, specInfoStruct)
		// pubCnt += 1
	}

	response(w, pubArr)
}

func getAosOrders(w http.ResponseWriter, r *http.Request) {
	phone := string(r.Form.Get("phone"))
	dateCurrent := string(r.Form.Get("date"))
	shopName := string(r.Form.Get("shop_name"))
	orderNo := string(r.Form.Get("order_no"))

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

	rows, err := orderInfoModel.GetList(shopName, phone, orderNo)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_aos.OrderInfoRecord, 0)

	//循环处理
	for rows.Next() {
		specInfoStruct, err := orderInfoModel.ScanRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, specInfoStruct)
	}

	response(w, pubArr)
}

func getSsdOrders(w http.ResponseWriter, r *http.Request) {
	phone := string(r.Form.Get("phone"))
	dateCurrent := string(r.Form.Get("date"))
	shopName := string(r.Form.Get("shop_name"))
	orderNo := string(r.Form.Get("order_no"))

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

	rows, err := orderInfoModel.GetList(shopName, phone, orderNo)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_order.OrderInfoRecord, 0)

	//循环处理
	for rows.Next() {
		specInfoStruct, err := orderInfoModel.ScanRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, specInfoStruct)
	}

	response(w, pubArr)
}

func getSkuSpecs(w http.ResponseWriter, r *http.Request) {
	shopName := string(r.Form.Get("shop_name"))
	dateCurrent := string(r.Form.Get("date"))

	//实例化数据模型
	specInfoModel, err := ssd_sku.CreateSpecInfo()
	serror.Check(err)
	defer func() {
		errClean := specInfoModel.CloseDB()
		serror.Check(errClean)
	}()

	err = specInfoModel.SetCurrentDate(dateCurrent)
	serror.Check(err)

	// log.Info("date:", dateCurrent)

	rows, err := specInfoModel.GetListByShopName(shopName, 500)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_sku.SpecInfoRecord, 0)

	//循环处理
	// pubCnt := 0
	for rows.Next() {
		specInfoStruct, err := specInfoModel.ScanRow(rows)
		serror.Check(err)

		//计入当前处理数组
		// pubArr[pubCnt] = specInfoStruct
		pubArr = append(pubArr, specInfoStruct)
		// pubCnt += 1
	}

	// if pubCnt != 0 {
	// 	pubArrSub := make([]ssd_sku.SpecInfoRecord, pubCnt)
	// 	copy(pubArrSub, pubArr)
	// 	//重新初始化数组切片
	// 	pubArr = make([]ssd_sku.SpecInfoRecord, 50)
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

	registerController("/get-sku-specs", getSkuSpecs)
	registerController("/get-sku-requests", getSkuRequests)
	registerController("/get-sku-responses", getSkuResponses)
	registerController("/get-ssd-orders", getSsdOrders)
	registerController("/get-aos-orders", getAosOrders)
	registerController("/get-tps-orders", getTpsOrders)
	registerController("/get-misc-devices", getMiscDevices)

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
