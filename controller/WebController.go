package controller

import (
	"encoding/json"
	"fmt"
	// "lib/config"
	// "lib/file"
	// "lib/log"
	"lib/response"
	"lib/serror"
	ssd_aos "model/aos"
	ssd_misc "model/misc"
	ssd_order "model/order"
	ssd_sku "model/sku"
	ssd_stat "model/stat"
	ssd_tps "model/tps"
	// "net"
	"io/ioutil"
	"net/http"
	"net/url"
	// "os"
	// "path/filepath"
	// "runtime"
)

// func httpDo() {
// 	client := &http.Client{}

// 	req, err := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))
// 	if err != nil {
// 		// handle error
// 	}

// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Set("Cookie", "name=anny")

// 	resp, err := client.Do(req)

// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		// handle error
// 	}

// 	fmt.Println(string(body))
// }

type resStruct struct {
	Body string
	Code int
}

func httpPostForm(addr string, vals map[string][]string) (resStruct, error) {
	values := url.Values{}
	for k, v := range vals {
		fmt.Println("key and value to post to remote: ", k, v)
		values[k] = v
	}
	// values := url.Values{"key": {"Value"}, "id": {"123"}}
	resp, err := http.PostForm(addr, values)
	if err != nil {
		// handle error
		return resStruct{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return resStruct{}, err
	}

	fmt.Println(string(body))

	bodyStr := string(body)

	res := resStruct{
		Body: bodyStr,
		Code: resp.StatusCode,
	}

	return res, nil
}

type WebController struct {
}

func (c *WebController) Login(w http.ResponseWriter, r *http.Request) {
	type resRecord struct {
		StatusCode int
		Msg        string
		Token      string `json:"token"`
	}
	rBody, _ := ioutil.ReadAll(r.Body) //把  body 内容读入字符串 rBody
	fmt.Println("body string from request: ", string(rBody))
	type loginRead struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	rBodyJson := loginRead{}
	err := json.Unmarshal(rBody, &rBodyJson)
	if err != nil {
		fmt.Println("error occurred when decode json, msg: ", err.Error())
		response.ResponseSuccess(w, &resRecord{
			StatusCode: 1,
			Msg:        err.Error(),
		})
	}
	username := string(rBodyJson.Username)
	password := string(rBodyJson.Password)
	fmt.Println("username from form:", username)
	fmt.Println("password from form:", password)

	vals := make(map[string][]string)
	vals["username"] = []string{username}
	vals["password"] = []string{password}

	_, err = httpPostForm("http://application-adm-api/auth/login", vals)
	if err != nil {
		response.ResponseError(w, 500)
	}

	res := &resRecord{
		Token: "ddddddddd",
	}

	response.ResponseSuccess(w, res)
}

func (c *WebController) GetSkuResponses(w http.ResponseWriter, r *http.Request) {
	shopNo := string(r.Form.Get("shop_no"))
	shopName := string(r.Form.Get("shop_name"))
	prodName := string(r.Form.Get("prod_name"))
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

	rows, err := goodsPurchaseModel.GetResponseList(searchNoArr, dateCurrent, prodName)
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

	response.ResponseSuccess(w, pubArr)
}

func (c *WebController) GetSkuRequests(w http.ResponseWriter, r *http.Request) {
	shopNo := string(r.Form.Get("shop_no"))
	shopName := string(r.Form.Get("shop_name"))
	prodName := string(r.Form.Get("prod_name"))
	//门店名称关键字对应的查询门店编码数组
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

	rows, err := goodsPurchaseModel.GetRequestList(searchNoArr, dateCurrent, prodName)
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

	response.ResponseSuccess(w, pubArr)
}

func (c *WebController) GetMiscDevices(w http.ResponseWriter, r *http.Request) {
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
				UpdateTime:  value.UpdateTime,
				AppVersion:  value.AppVersion,
				SerialNo:    value.SerialNo,
				NetworkType: value.NetworkType,
				AppEnv:      value.AppEnv,
				IsCanary:    isCanary,
				CreateTime:  value.CreateTime,
			})
		}
	}

	response.ResponseSuccess(w, resArr)
}

func (c *WebController) GetTpsOrders(w http.ResponseWriter, r *http.Request) {
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

	response.ResponseSuccess(w, pubArr)
}

func (c *WebController) GetAosOrders(w http.ResponseWriter, r *http.Request) {
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

	response.ResponseSuccess(w, pubArr)
}

func (c *WebController) GetSsdOrders(w http.ResponseWriter, r *http.Request) {
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

	response.ResponseSuccess(w, pubArr)
}

func (c *WebController) GetSkuSpecs(w http.ResponseWriter, r *http.Request) {
	shopName := string(r.Form.Get("shop_name"))
	shopCode := string(r.Form.Get("shop_code"))
	prodName := string(r.Form.Get("prod_name"))
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

	rows, err := specInfoModel.GetList(shopCode, shopName, prodName, 500)
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
	response.ResponseSuccess(w, pubArr)
	// responseError(w, 401)
}
