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
	ID          int    `json:"id"`
	ShopNo      string `json:"shop_no"`
	AppVersion  string `json:"app_version"`
	SerialNo    string `json:"serial_no"`
	NetworkType string `json:"network_type"`
	AppEnv      string `json:"app_env"`
	UpdateTime  string `json:"update_time"`
	CreateTime  string `json:"create_time"`
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

//GetList ...
func (oi *Device) GetList(shopNoArr []string, appVersion string) (*sql.Rows, error) {
	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	//构造条件语句
	var whereArr []string

	where := ""

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

	limit := "5000"
	fields := oi.getFields()

	var rows *sql.Rows
	var err error
	if where != "" {
		rows, err = oi.prepare().Select(fields).Where(where).Order("update_time desc").Limit(limit).Rows()
	} else {
		rows, err = oi.prepare().Select(fields).Order("update_time desc").Limit(limit).Rows()
	}

	if err != nil {
		return nil, err
	}
	return rows, nil
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
