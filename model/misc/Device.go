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
	"strings"
	// "strconv"
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

func (oi *Device) getFields() []string {
	var keys []string
	var oRecord Record
	dataType := reflect.TypeOf(oRecord)
	for k := 0; k < dataType.NumField(); k++ {
		keys = append(keys, dataType.Field(k).Tag.Get("json"))
	}
	return keys
}

func (oi *Device) buildConditions(shopNoArr []string, appVersion string, companySaleId string, model string) string {
	where := ""
	//构造条件语句
	var whereArr []string

	whereArr = append(whereArr, fmt.Sprintf("is_deleted=%d", 0))
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
		whereArr = append(whereArr, fmt.Sprintf("shop_no in (%s)", shopNoStr))
	}
	if appVersion != "" {
		whereArr = append(whereArr, fmt.Sprintf("app_version='%s'", appVersion))
	}
	if companySaleId != "" {
		whereArr = append(whereArr, fmt.Sprintf("company_sale_id='%s'", companySaleId))
	}
	if model != "" {
		whereArr = append(whereArr, fmt.Sprintf("model='%s'", model))
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

	return where
}

//GetList ...
func (oi *Device) GetList(shopNoArr []string, appVersion string, companySaleId string, model string, pageSize int, pageNum int) (*sql.Rows, error) {
	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	where := oi.buildConditions(shopNoArr, appVersion, companySaleId, model)
	limit := pageSize
	offset := pageSize * (pageNum - 1)
	fields := oi.getFields()

	var rows *sql.Rows
	var err error
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
func (oi *Device) GetComplexList(shopNoArr []string, appVersion string, companySaleId string, model string, isCanary int, pageSize int, pageNum int) (*sql.Rows, error) {
	where := ""
	//构造条件语句
	var whereArr []string

	whereArr = append(whereArr, fmt.Sprintf("device.is_deleted=%d", 0))
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
		whereArr = append(whereArr, fmt.Sprintf("device.shop_no in (%s)", shopNoStr))
	}
	if appVersion != "" {
		whereArr = append(whereArr, fmt.Sprintf("device.app_version='%s'", appVersion))
	}
	if companySaleId != "" {
		whereArr = append(whereArr, fmt.Sprintf("device.company_sale_id='%s'", companySaleId))
	}
	if model != "" {
		whereArr = append(whereArr, fmt.Sprintf("device.model='%s'", model))
	}
	// 灰度设备
	if isCanary == 1 {
		whereArr = append(whereArr, fmt.Sprintf("canary_device.id<>0"))
	} else if isCanary == 0 {
		// 非灰度设备
		whereArr = append(whereArr, fmt.Sprintf("canary_device.id IS NULL"))
	}
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
	limit := pageSize
	offset := pageSize * (pageNum - 1)
	fields := ""
	// var keys []string
	var oRecord Record
	dataType := reflect.TypeOf(oRecord)
	for k := 0; k < dataType.NumField(); k++ {
		if k == 0 {
			fields = fmt.Sprintf("device.%s", dataType.Field(k).Tag.Get("json"))
		} else {
			fields = fmt.Sprintf("%s, device.%s", fields, dataType.Field(k).Tag.Get("json"))
		}
	}

	fields = fmt.Sprintf("%s, canary_device.id as canary_id", fields)

	join := "left join canary_device on device.id=canary_device.device_id"
	order := "device.id desc"

	var rows *sql.Rows
	var err error
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
func (oi *Device) GetComplexCount(shopNoArr []string, appVersion string, companySaleId string, model string, isCanary int) int {
	where := ""
	//构造条件语句
	var whereArr []string

	whereArr = append(whereArr, fmt.Sprintf("device.is_deleted=%d", 0))
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
		whereArr = append(whereArr, fmt.Sprintf("device.shop_no in (%s)", shopNoStr))
	}
	if appVersion != "" {
		whereArr = append(whereArr, fmt.Sprintf("device.app_version='%s'", appVersion))
	}
	if companySaleId != "" {
		whereArr = append(whereArr, fmt.Sprintf("device.company_sale_id='%s'", companySaleId))
	}
	if model != "" {
		whereArr = append(whereArr, fmt.Sprintf("device.model='%s'", model))
	}
	// 灰度设备
	if isCanary == 1 {
		whereArr = append(whereArr, fmt.Sprintf("canary_device.id<>0"))
	} else if isCanary == 0 {
		// 非灰度设备
		whereArr = append(whereArr, fmt.Sprintf("canary_device.id IS NULL"))
	}
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

	fields := "count(*)"
	join := "left join canary_device on device.id=canary_device.device_id"
	order := "device.id desc"

	count := 0

	if where != "" {
		oi.prepare().Select(fields).Joins(join).Where(where).Order(order).Count(&count)
	} else {
		oi.prepare().Select(fields).Joins(join).Order(order).Count(&count)
	}

	return count
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
func (oi *Device) GetCount(shopNoArr []string, appVersion string, companySaleId string, model string) int {
	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	// alias := "d"
	where := oi.buildConditions(shopNoArr, appVersion, companySaleId, model)
	fields := "count(*)"
	count := 0

	if where != "" {
		oi.prepare().Select(fields).Where(where).Count(&count)
	} else {
		oi.prepare().Select(fields).Count(&count)
	}

	return count
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
