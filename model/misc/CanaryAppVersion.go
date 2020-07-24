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

//CanaryAppVersionRecord ...
//数据表字段结构
type CanaryAppVersionRecord struct {
	ID         int    `json:"id"`
	IsForce    int    `json:"is_force"`
	Version    string `json:"version"`
	Url        string `json:"url"`
	CreateTime int    `json:"create_time"`
}

//CanaryAppVersion ... 数据模型对象
type CanaryAppVersion struct {
	model.Model
}

func (oi *CanaryAppVersion) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func (oi *CanaryAppVersion) prepare() *gorm.DB {
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

func (oi *CanaryAppVersion) getFields() []string {
	var keys []string
	var oRecord CanaryAppVersionRecord
	dataType := reflect.TypeOf(oRecord)
	for k := 0; k < dataType.NumField(); k++ {
		keys = append(keys, dataType.Field(k).Tag.Get("json"))
	}
	return keys
}

//GetList ...
func (oi *CanaryAppVersion) GetList() (*sql.Rows, error) {
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
func (oi *CanaryAppVersion) ScanRow(r *sql.Rows) (CanaryAppVersionRecord, error) {
	var oRecord CanaryAppVersionRecord
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

//CreateCanaryAppVersion ...
func CreateCanaryAppVersion() (*CanaryAppVersion, error) {
	obj := &CanaryAppVersion{}
	err := obj.OpenDB("common.db.ssd_misc")
	if err != nil {
		return nil, err
	}
	err = obj.SetTableName("canary_app_version")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
