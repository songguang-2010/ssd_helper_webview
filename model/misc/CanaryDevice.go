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
	// "strings"
	// "strconv"
)

//CanaryDeviceRecord ...
//数据表字段结构
type CanaryDeviceRecord struct {
	ID         int    `json:"id"`
	ShopNo     string `json:"shop_no"`
	AppVersion string `json:"app_version"`
	DeviceID   int    `json:"device_id"`
}

//CanaryDevice ... 数据模型对象
type CanaryDevice struct {
	model.Model
}

func (oi *CanaryDevice) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func (oi *CanaryDevice) prepare() *gorm.DB {
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

func (oi *CanaryDevice) getFields() []string {
	var keys []string
	var oRecord CanaryDeviceRecord
	dataType := reflect.TypeOf(oRecord)
	for k := 0; k < dataType.NumField(); k++ {
		keys = append(keys, dataType.Field(k).Tag.Get("json"))
	}
	return keys
}

//GetList ...
func (oi *CanaryDevice) GetList() (*sql.Rows, error) {
	//构造条件语句
	var whereArr []string

	where := ""

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

	fields := oi.getFields()

	var rows *sql.Rows
	var err error
	if where != "" {
		rows, err = oi.prepare().Select(fields).Where(where).Order("id asc").Rows()
	} else {
		rows, err = oi.prepare().Select(fields).Order("id asc").Rows()
	}

	if err != nil {
		return nil, err
	}
	return rows, nil
}

//ScanRow ...
func (oi *CanaryDevice) ScanRow(r *sql.Rows) (CanaryDeviceRecord, error) {
	var oRecord CanaryDeviceRecord
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

//CreateCanaryDevice ...
func CreateCanaryDevice() (*CanaryDevice, error) {
	obj := &CanaryDevice{}
	err := obj.OpenDB("common.db.ssd_misc")
	if err != nil {
		return nil, err
	}
	err = obj.SetTableName("canary_device")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
