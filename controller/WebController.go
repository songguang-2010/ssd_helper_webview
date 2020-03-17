package controller

import (
	"encoding/json"
	"fmt"
	// "reflect"
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
	// service_ssd_misc "service/misc"
	// "net"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	// "strings"
	"time"
	// "os"
	// "path/filepath"
	// "runtime"
	"strconv"
)

var uid int
var token string

//生成随机字符串
func GetRandomString(randomLen int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomLen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(8)
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 生成秘钥
func createSecret(uid int) string {
	nonceStrPool := "asdfghjkl0987oiuytrewq654zxcvbnm321"
	step := 5
	startStrOffset := uid % step
	signStr := ""
	for i := 0; i < 35; i += step {
		strOffset := startStrOffset + i
		userIdTmp := uid + i
		signStr = fmt.Sprintf("%s%s%d", signStr, nonceStrPool[strOffset:strOffset+1], userIdTmp)
	}
	ctx := md5.New()
	ctx.Write([]byte(signStr))
	return hex.EncodeToString(ctx.Sum(nil))
}

func createSignature(token string, uid int, timestamp int64, nonce string, values url.Values) string {
	secret := createSecret(uid)
	postBody := values.Encode()
	signStr := fmt.Sprintf("nonce=%s&timestamp=%d&token=%s&uid=%d&postBody=%s%s", nonce, timestamp, token, uid, postBody, secret)
	ctx := md5.New()
	ctx.Write([]byte(signStr))
	return hex.EncodeToString(ctx.Sum(nil))
}

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

// type resStruct struct {
// 	Body string
// 	Code int
// }

func httpPostForm(addr string, vals map[string][]string) (interface{}, error) {
	// 返回值结构
	var resStruct interface{}

	values := url.Values{}
	for k, v := range vals {
		fmt.Println("key and value to post to remote: ", k, v[0])
		values[k] = v
	}
	// values := url.Values{"key": {"Value"}, "id": {"123"}}

	//通用验证参数
	timestampNow := time.Now().Unix()
	nonceStr := GetRandomSalt()
	signature := createSignature(token, uid, timestampNow, nonceStr, values)
	getParams := fmt.Sprintf("token=%s&uid=%d&timestamp=%d&nonce=%s&signature=%s", token, uid, timestampNow, nonceStr, signature)
	postAddr := fmt.Sprintf("http://172.31.1.21:8101%s?%s", addr, getParams)
	fmt.Println(postAddr)
	resp, err := http.PostForm(postAddr, values)
	if err != nil {
		// handle error
		fmt.Println(err.Error())
		return resStruct, err
	}
	if resp.StatusCode != 200 {
		// handle error
		errMsg := fmt.Sprintf("error occurred, error code: %d, error msg: %s", resp.StatusCode, resp.Status)
		return resStruct, serror.New(errMsg)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return resStruct, err
	}

	fmt.Println(string(body))

	// 远端响应值结构
	type receiverStruct struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	rBodyJson := receiverStruct{}
	err = json.Unmarshal(body, &rBodyJson)
	if err != nil {
		fmt.Println("error occurred when decode json, msg: ", err.Error())
		// handle error
		return resStruct, err
	}

	code := int(rBodyJson.Code)
	if code != 100 {
		return resStruct, serror.New(rBodyJson.Msg)
	}

	return rBodyJson.Data, nil
}

type WebController struct {
}

// SetCanaryBatch 批量设置为灰度设备
func (c *WebController) SetCanaryBatch(w http.ResponseWriter, r *http.Request) {

	// 接收客户端请求值结构
	type paramsClient struct {
		DeviceIds string `json:"device_ids"`
	}

	// 读取客户端请求值
	rBody, _ := ioutil.ReadAll(r.Body)
	rBodyJson := paramsClient{}
	err := json.Unmarshal(rBody, &rBodyJson)
	if err != nil {
		fmt.Println("error occurred when decode json, msg: ", err.Error())
		response.ResponseSuccess(w, 500)
	}
	device_ids := fmt.Sprintf("%s", rBodyJson.DeviceIds)

	fmt.Println("device ids from form: ", device_ids)

	vals := make(map[string][]string)
	vals["deviceIds"] = []string{device_ids}

	resRemote, err := httpPostForm("/device/set-canary-batch", vals)
	if err != nil {
		fmt.Println("error occurred when response from remote, msg: ", err.Error())
		response.ResponseError(w, 500)
	}

	fmt.Printf("response from remote: %v", resRemote)

	resRecord := make(map[string]string)
	resRecord["result"] = "ok"

	response.ResponseSuccess(w, resRecord)
}

// CancelCanaryBatch 批量取消为灰度设备
func (c *WebController) CancelCanaryBatch(w http.ResponseWriter, r *http.Request) {

	// 接收客户端请求值结构
	type paramsClient struct {
		DeviceIds string `json:"device_ids"`
	}

	// 读取客户端请求值
	rBody, _ := ioutil.ReadAll(r.Body)
	rBodyJson := paramsClient{}
	err := json.Unmarshal(rBody, &rBodyJson)
	if err != nil {
		fmt.Println("error occurred when decode json, msg: ", err.Error())
		response.ResponseSuccess(w, 500)
	}
	device_ids := fmt.Sprintf("%s", rBodyJson.DeviceIds)

	fmt.Println("device ids from form: ", device_ids)

	vals := make(map[string][]string)
	vals["deviceIds"] = []string{device_ids}

	resRemote, err := httpPostForm("/device/cancel-canary-batch", vals)
	if err != nil {
		fmt.Println("error occurred when response from remote, msg: ", err.Error())
		response.ResponseError(w, 500)
	}

	fmt.Printf("response from remote: %v", resRemote)

	resRecord := make(map[string]string)
	resRecord["result"] = "ok"

	response.ResponseSuccess(w, resRecord)
}

// CancelCanary 取消为灰度设备
func (c *WebController) CancelCanary(w http.ResponseWriter, r *http.Request) {

	// 接收客户端请求值结构
	type paramsClient struct {
		DeviceId int `json:"device_id"`
	}

	// 读取客户端请求值
	rBody, _ := ioutil.ReadAll(r.Body)
	rBodyJson := paramsClient{}
	err := json.Unmarshal(rBody, &rBodyJson)
	if err != nil {
		fmt.Println("error occurred when decode json, msg: ", err.Error())
		response.ResponseSuccess(w, 500)
	}
	device_id := fmt.Sprintf("%d", rBodyJson.DeviceId)

	fmt.Println("device id from form: ", device_id)

	vals := make(map[string][]string)
	vals["deviceId"] = []string{device_id}

	resRemote, err := httpPostForm("/device/cancel-canary", vals)
	if err != nil {
		fmt.Println("error occurred when response from remote, msg: ", err.Error())
		response.ResponseError(w, 500)
	}

	fmt.Printf("response from remote: %v", resRemote)

	resRecord := make(map[string]string)
	resRecord["result"] = "ok"

	response.ResponseSuccess(w, resRecord)
}

// SetCanary 设置为灰度设备
func (c *WebController) SetCanary(w http.ResponseWriter, r *http.Request) {

	// 接收客户端请求值结构
	type paramsClient struct {
		DeviceId int `json:"device_id"`
	}

	// 读取客户端请求值
	rBody, _ := ioutil.ReadAll(r.Body)
	rBodyJson := paramsClient{}
	err := json.Unmarshal(rBody, &rBodyJson)
	if err != nil {
		fmt.Println("error occurred when decode json, msg: ", err.Error())
		response.ResponseSuccess(w, 500)
	}
	device_id := fmt.Sprintf("%d", rBodyJson.DeviceId)

	fmt.Println("device id from form: ", device_id)

	vals := make(map[string][]string)
	vals["deviceId"] = []string{device_id}

	resRemote, err := httpPostForm("/device/set-canary", vals)
	if err != nil {
		fmt.Println("error occurred when response from remote, msg: ", err.Error())
		response.ResponseError(w, 500)
	}

	fmt.Printf("response from remote: %v", resRemote)

	resRecord := make(map[string]string)
	resRecord["result"] = "ok"

	response.ResponseSuccess(w, resRecord)
}

func (c *WebController) Login(w http.ResponseWriter, r *http.Request) {

	// 接收客户端请求值结构
	type loginRead struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 读取客户端请求值
	rBody, _ := ioutil.ReadAll(r.Body)
	// fmt.Println("body string from request: ", string(rBody))
	rBodyJson := loginRead{}
	err := json.Unmarshal(rBody, &rBodyJson)
	if err != nil {
		fmt.Println("error occurred when decode json, msg: ", err.Error())
		response.ResponseSuccess(w, 500)
	}
	username := string(rBodyJson.Username)
	password := string(rBodyJson.Password)
	fmt.Println("username from form:", username)
	fmt.Println("password from form:", password)

	vals := make(map[string][]string)
	vals["username"] = []string{username}
	vals["password"] = []string{password}

	resRemote, err := httpPostForm("/auth/login", vals)
	if err != nil {
		response.ResponseError(w, 500)
	}

	fmt.Printf("response from remote: %v", resRemote)

	// type resStruct
	resMap := resRemote.(map[string]interface{})

	uid = (int)((resMap["uid"]).(float64))
	token = resMap["token"].(string)

	fmt.Printf("response uid from remote: %d", uid)
	fmt.Printf("response token from remote: %s", token)

	// //返回值结构
	// type resRecord struct {
	// 	StatusCode int
	// 	Msg        string
	// 	Token      string `json:"token"`
	// }

	// res := &resRecord{
	// 	Token: "ddddddddd",
	// }

	resRecord := make(map[string]string)
	resRecord["token"] = token

	response.ResponseSuccess(w, resRecord)
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
	//接收客户端参数
	shopNo := string(r.Form.Get("shop_no"))
	shopName := string(r.Form.Get("shop_name"))
	appVersion := string(r.Form.Get("app_version"))
	companySaleId := string(r.Form.Get("company_sale_id"))
	model := string(r.Form.Get("model"))

	// 是否灰度设备
	isCanary, err := strconv.Atoi(string(r.Form.Get("is_canary")))
	serror.Check(err)

	//分页参数
	pageNum, err := strconv.Atoi(string(r.Form.Get("pageNum")))
	serror.Check(err)
	pageSize, err := strconv.Atoi(string(r.Form.Get("pageSize")))
	serror.Check(err)
	fmt.Println("current page num: ", pageNum)
	fmt.Println("current page size: ", pageSize)
	// canary := string(r.Form.Get("canary"))

	//声明门店编码过滤数组
	searchShopNoArr := make([]string, 0)

	// isCanary, err := strconv.Atoi(canary)
	// serror.Check(err)

	//实例化门店信息缓存数据模型
	shopModel, err := ssd_stat.CreateJwdShopCache()
	serror.Check(err)
	defer func() {
		errClean := shopModel.CloseDB()
		serror.Check(errClean)
	}()

	//如果搜索参数中含有门店名称
	if shopName != "" {
		//根据门店名称查询相关门店列表
		shopList, err := shopModel.GetListByName(shopName)
		serror.Check(err)
		defer shopList.Close()

		//循环处理, 提取门店编码
		for shopList.Next() {
			shopInfo, err := shopModel.ScanRow(shopList)
			serror.Check(err)
			//计入当前门店编码过滤数组
			searchShopNoArr = append(searchShopNoArr, shopInfo.ShopNo)
		}
	}

	//如果搜索参数中含有门店编码
	if shopNo != "" {
		//将门店编码计入当前门店编码过滤数组
		searchShopNoArr = append(searchShopNoArr, shopNo)
	}

	//实例化数据模型
	deviceModel, err := ssd_misc.CreateDevice()
	serror.Check(err)
	defer func() {
		errClean := deviceModel.CloseDB()
		serror.Check(errClean)
	}()

	// 获取设备数量
	count := deviceModel.GetComplexCount(searchShopNoArr, appVersion, companySaleId, model, isCanary)
	// 获取设备列表数据
	rows, err := deviceModel.GetComplexList(searchShopNoArr, appVersion, companySaleId, model, isCanary, pageSize, pageNum)
	serror.Check(err)
	defer rows.Close()

	//初始化数组
	pubArr := make([]ssd_misc.ComplexRecord, 0)

	//循环处理
	for rows.Next() {
		oRecord, err := deviceModel.ScanComplexRow(rows)
		serror.Check(err)
		//计入当前处理数组
		pubArr = append(pubArr, oRecord)
	}

	//初始化门店编码及名称映射列表
	shopMap := make(map[string]string)
	if len(pubArr) != 0 {
		//获取门店编码列表
		shopNoArr := make([]string, 0)
		for _, value := range pubArr {
			shopNoArr = append(shopNoArr, value.ShopNo)
		}
		//获取门店信息列表
		rowsShopList, err := shopModel.GetListByShopNoArr(shopNoArr)
		serror.Check(err)
		defer rowsShopList.Close()
		//遍历读取门店信息
		for rowsShopList.Next() {
			oRecord, err := shopModel.ScanRow(rowsShopList)
			serror.Check(err)
			//计入当前门店编码及名称映射列表
			shopMap[oRecord.ShopNo] = oRecord.ShopName
		}
	}

	//初始化灰度设备id映射列表
	// deviceIDMap := make(map[int]int)
	//实例化灰度设备服务模型
	// canaryDeviceService, err := service_ssd_misc.CreateCanaryDevice()
	// serror.Check(err)
	//查询灰度设备id列表
	// deviceIDArr, err := canaryDeviceService.GetDeviceIdList()
	// serror.Check(err)
	// for _, v := range deviceIDArr {
	//计入当前灰度设备列表
	// deviceIDMap[v] = 1
	// }

	//定义返回数据类型
	type resRecord struct {
		ID            int    `json:"id"`
		ShopNo        string `json:"shop_no"`
		ShopName      string `json:"shop_name"`
		AppVersion    string `json:"app_version"`
		SerialNo      string `json:"serial_no"`
		NetworkType   string `json:"network_type"`
		AppEnv        string `json:"app_env"`
		IsCanary      string `json:"is_canary"`
		UpdateTime    string `json:"update_time"`
		CreateTime    string `json:"create_time"`
		CompanySaleId string `json:"company_sale_id"`
	}

	//初始化返回数组
	resArr := make([]resRecord, 0)

	if len(pubArr) != 0 {
		// 构造返回数组
		for _, value := range pubArr {
			//是否灰度设备
			isCanary := "0"
			// if deviceIDMap[value.ID] != 0 {
			// isCanary = "1"
			// }
			if value.CanaryId != 0 {
				isCanary = "1"
			}

			// 返回数组
			resArr = append(resArr, resRecord{
				ID:            value.ID,
				ShopNo:        value.ShopNo,
				ShopName:      shopMap[value.ShopNo],
				UpdateTime:    value.UpdateTime,
				AppVersion:    value.AppVersion,
				SerialNo:      value.SerialNo,
				NetworkType:   value.NetworkType,
				AppEnv:        value.AppEnv,
				IsCanary:      isCanary,
				CreateTime:    value.CreateTime,
				CompanySaleId: value.CompanySaleId,
			})
		}
	}

	//构造实际返回数据
	type resData struct {
		Total int         `json:"total"`
		List  []resRecord `json:"list"`
	}

	response.ResponseSuccess(w, resData{
		Total: count,
		List:  resArr,
	})
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
