package misc

import (
	"database/sql"
	"fmt"
	// "github.com/huandu/xstrings"
	"github.com/jinzhu/gorm"
	//just init
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"lib/model"
	// "lib/stime"
	// "crypto/md5"
	// "github.com/huandu/xstrings"
	// "io"
	"reflect"
	"strconv"
	"strings"
)

//Record ...
//数据表字段结构
type Record struct {
	ID            int    `json:"id"`
	ShopNo        string `json:"shop_no"`
	AppVersion    string `json:"app_version"`
	SerialNo      string `json:"serial_no"`
	NetworkType   string `json:"network_type"`
	AppEnv        string `json:"app_env"`
	UpdateTime    string `json:"update_time"`
	CreateTime    string `json:"create_time"`
	CompanySaleId string `json:"company_sale_id"`
}

type ComplexRecord struct {
	ID            int    `json:"id"`
	ShopNo        string `json:"shop_no"`
	AppVersion    string `json:"app_version"`
	SerialNo      string `json:"serial_no"`
	NetworkType   string `json:"network_type"`
	AppEnv        string `json:"app_env"`
	UpdateTime    string `json:"update_time"`
	CreateTime    string `json:"create_time"`
	CompanySaleId string `json:"company_sale_id"`
	CanaryId      int    `json:"canary_id"`
}

//Device ... 数据模型对象
type Device struct {
	model.Model
}

func (oi *Device) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func (oi *Device) prepare() *gorm.DB {
	table, err := oi.getTableName()
	if err != nil {
		oi.SetError(err.Error())
		return nil
	}
	db, err := oi.GetDB()
	if err != nil {
		oi.SetError(err.Error())
		return nil
	}
	return db.Table(table)
}

func (oi *Device) getFields() (string, error) {
	table, err := oi.getTableName()
	if err != nil {
		return "", err
	}

	fields := ""

	var oRecord Record
	dataType := reflect.TypeOf(oRecord)
	for k := 0; k < dataType.NumField(); k++ {
		if k == 0 {
			fields = fmt.Sprintf("%s.%s", table, dataType.Field(k).Tag.Get("json"))
		} else {
			fields = fmt.Sprintf("%s, %s.%s", fields, table, dataType.Field(k).Tag.Get("json"))
		}
	}

	return fields, nil
}

func (oi *Device) buildConditions(shopNoArr []string, appVersion string, companySaleId string, model string, appEnv string) (string, error) {
	table, err := oi.getTableName()
	if err != nil {
		return "", err
	}

	where := ""
	//构造条件语句
	var whereArr []string

	whereArr = append(whereArr, fmt.Sprintf("%s.is_deleted=%d", table, 0))
	if len(shopNoArr) != 0 {
		shopNoStr := ""
		for k, v := range shopNoArr {
			if k == 0 {
				shopNoStr = fmt.Sprintf("'%s'", strings.TrimSpace(v))
			} else {
				shopNoStr = fmt.Sprintf("%s,'%s'", shopNoStr, strings.TrimSpace(v))
			}
		}
		fmt.Println(shopNoStr)
		whereArr = append(whereArr, fmt.Sprintf("%s.shop_no in (%s)", table, shopNoStr))
	}
	if appVersion != "" {
		whereArr = append(whereArr, fmt.Sprintf("%s.app_version='%s'", table, appVersion))
	}
	if companySaleId != "" {
		whereArr = append(whereArr, fmt.Sprintf("%s.company_sale_id='%s'", table, companySaleId))
	}
	if model != "" {
		whereArr = append(whereArr, fmt.Sprintf("%s.model='%s'", table, model))
	}
	if appEnv != "" {
		whereArr = append(whereArr, fmt.Sprintf("%s.app_env='%s'", table, appEnv))
	}
	// if len(deviceIDArr) != 0 {
	// 	deviceIDStr := ""
	// 	for k, v := range deviceIDArr {
	// 		if k == 0 {
	// 			deviceIDStr = fmt.Sprintf("'%s'", strings.TrimSpace(v))
	// 		} else {
	// 			deviceIDStr = fmt.Sprintf("%s,'%s'", deviceIDStr, strings.TrimSpace(v))
	// 		}
	// 	}
	// 	// fmt.Println(deviceIDStr)
	// 	whereArr = append(whereArr, fmt.Sprintf("id in (%s)", deviceIDStr))
	// }

	if len(whereArr) > 0 {
		for k, v := range whereArr {
			// fmt.Println(k, v)
			if k == 0 {
				where = fmt.Sprintf("%s", v)
			} else {
				where = fmt.Sprintf("%s and %s", where, v)
			}
		}
	}

	return where, nil
}

//GetList ...
func (oi *Device) GetList(shopNoArr []string, paramsMap map[string]string) (*sql.Rows, error) {
	appVersion := paramsMap["appVersion"]
	companySaleId := paramsMap["companySaleId"]
	model := paramsMap["model"]
	appEnv := paramsMap["appEnv"]
	// 如果不过滤灰度状态，则设置为空
	if appEnv == "all" {
		appEnv = ""
	}
	// isCanary, err := strconv.Atoi(paramsMap["isCanary"])
	// if err != nil {
	// 	return nil, err
	// }
	// testExclude, err := strconv.Atoi(paramsMap["testExclude"])
	// if err != nil {
	// 	return nil, err
	// }
	pageSize, err := strconv.Atoi(paramsMap["pageSize"])
	if err != nil {
		return nil, err
	}
	pageNum, err := strconv.Atoi(paramsMap["pageNum"])
	if err != nil {
		return nil, err
	}
	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	where, err := oi.buildConditions(shopNoArr, appVersion, companySaleId, model, appEnv)
	if err != nil {
		return nil, err
	}
	limit := pageSize
	offset := pageSize * (pageNum - 1)
	fields, err := oi.getFields()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows

	if where != "" {
		rows, err = oi.prepare().Select(fields).Where(where).Order("id desc").Offset(offset).Limit(limit).Rows()
	} else {
		rows, err = oi.prepare().Select(fields).Order("id desc").Offset(offset).Limit(limit).Rows()
	}

	if err != nil {
		return nil, err
	}
	return rows, nil
}

//GetComplexList ...
func (oi *Device) GetComplexList(shopNoArr []string, paramsMap map[string]string) (*sql.Rows, error) {
	appVersion := paramsMap["appVersion"]
	companySaleId := paramsMap["companySaleId"]
	model := paramsMap["model"]
	appEnv := paramsMap["appEnv"]
	// 如果不过滤灰度状态，则设置为空
	if appEnv == "all" {
		appEnv = ""
	}
	isCanary, err := strconv.Atoi(paramsMap["isCanary"])
	if err != nil {
		return nil, err
	}
	testExclude, err := strconv.Atoi(paramsMap["testExclude"])
	if err != nil {
		return nil, err
	}
	pageSize, err := strconv.Atoi(paramsMap["pageSize"])
	if err != nil {
		return nil, err
	}
	pageNum, err := strconv.Atoi(paramsMap["pageNum"])
	if err != nil {
		return nil, err
	}

	where, err := oi.buildConditions(shopNoArr, appVersion, companySaleId, model, appEnv)
	if err != nil {
		return nil, err
	}

	// 排除测试设备
	if testExclude == 1 {
		where = fmt.Sprintf("%s and device.shop_no<>'000001'", where)
	}

	// 灰度设备
	if isCanary == 1 {
		where = fmt.Sprintf("%s and canary_device.id<>0", where)
	} else if isCanary == 0 {
		// 非灰度设备
		where = fmt.Sprintf("%s and canary_device.id IS NULL", where)
	}

	limit := pageSize
	offset := pageSize * (pageNum - 1)

	fields, err := oi.getFields()
	if err != nil {
		return nil, err
	}

	fields = fmt.Sprintf("%s, canary_device.id as canary_id", fields)

	join := "left join canary_device on device.id=canary_device.device_id"
	order := "device.id desc"

	var rows *sql.Rows

	if where != "" {
		rows, err = oi.prepare().Select(fields).Joins(join).Where(where).Order(order).Offset(offset).Limit(limit).Rows()
	} else {
		rows, err = oi.prepare().Select(fields).Joins(join).Order(order).Offset(offset).Limit(limit).Rows()
	}

	if err != nil {
		return nil, err
	}

	return rows, nil
}

//GetComplexCount ...
func (oi *Device) GetComplexCount(shopNoArr []string, paramsMap map[string]string) (int, error) {
	appVersion := paramsMap["appVersion"]
	companySaleId := paramsMap["companySaleId"]
	model := paramsMap["model"]
	appEnv := paramsMap["appEnv"]
	// 如果不过滤灰度状态，则设置为空
	if appEnv == "all" {
		appEnv = ""
	}
	isCanary, err := strconv.Atoi(paramsMap["isCanary"])
	if err != nil {
		return 0, err
	}
	testExclude, err := strconv.Atoi(paramsMap["testExclude"])
	if err != nil {
		return 0, err
	}

	where, err := oi.buildConditions(shopNoArr, appVersion, companySaleId, model, appEnv)
	if err != nil {
		return 0, err
	}

	// 排除测试设备
	if testExclude == 1 {
		where = fmt.Sprintf("%s and device.shop_no<>'000001'", where)
	}

	// 灰度设备
	if isCanary == 1 {
		where = fmt.Sprintf("%s and canary_device.id<>0", where)
	} else if isCanary == 0 {
		// 非灰度设备
		where = fmt.Sprintf("%s and canary_device.id IS NULL", where)
	}

	fields := "count(*)"
	join := "left join canary_device on device.id=canary_device.device_id"
	order := "device.id desc"

	count := 0

	if where != "" {
		oi.prepare().Select(fields).Joins(join).Where(where).Order(order).Count(&count)
	} else {
		oi.prepare().Select(fields).Joins(join).Order(order).Count(&count)
	}

	return count, nil
}

//ScanComplexRow ...
func (oi *Device) ScanComplexRow(r *sql.Rows) (ComplexRecord, error) {
	var oRecord ComplexRecord
	db, err := oi.GetDB()
	if err != nil {
		// oi.SetError(err.Error())
		return oRecord, err
	}
	err = db.ScanRows(r, &oRecord)
	if err != nil {
		return oRecord, err
	}
	return oRecord, nil
}

//GetCount ...
func (oi *Device) GetCount(shopNoArr []string, paramsMap map[string]string) (int, error) {
	appVersion := paramsMap["appVersion"]
	companySaleId := paramsMap["companySaleId"]
	model := paramsMap["model"]
	appEnv := paramsMap["appEnv"]
	// 如果不过滤灰度状态，则设置为空
	if appEnv == "all" {
		appEnv = ""
	}
	// isCanary, err := strconv.Atoi(paramsMap["isCanary"])
	// if err != nil {
	// 	return 0, err
	// }
	// testExclude, err := strconv.Atoi(paramsMap["testExclude"])
	// if err != nil {
	// 	return 0, err
	// }
	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	// alias := "d"
	where, err := oi.buildConditions(shopNoArr, appVersion, companySaleId, model, appEnv)
	if err != nil {
		return 0, err
	}
	fields := "count(*)"
	count := 0

	if where != "" {
		oi.prepare().Select(fields).Where(where).Count(&count)
	} else {
		oi.prepare().Select(fields).Count(&count)
	}

	return count, nil
}

//ScanRow ...
func (oi *Device) ScanRow(r *sql.Rows) (Record, error) {
	var oRecord Record
	db, err := oi.GetDB()
	if err != nil {
		// oi.SetError(err.Error())
		return oRecord, err
	}
	err = db.ScanRows(r, &oRecord)
	if err != nil {
		return oRecord, err
	}
	return oRecord, nil
}

//CreateDevice ...
func CreateDevice() (*Device, error) {
	obj := &Device{}
	err := obj.OpenDB("common.db.ssd_misc")
	if err != nil {
		return nil, err
	}
	err = obj.SetTableName("device")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
