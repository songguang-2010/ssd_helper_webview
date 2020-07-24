package aos

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
	ID          int    `json:"id"`
	Type        int    `json:"type"`
	Params      string `json:"params"`
	Status      int    `json:"status"`
	Result      string `json:"result"`
	Description string `json:"description"`
	UpdateTime  string `json:"update_time"`
	CreateTime  string `json:"create_time"`
}

//数据模型对象
type Job struct {
	model.Model
}

func (oi *Job) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func (oi *Job) prepare() *gorm.DB {
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

func (oi *Job) getFields() (string, error) {
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

func (oi *Job) buildConditions(jobType int, status int) (string, error) {
	table, err := oi.getTableName()
	if err != nil {
		return "", err
	}

	where := ""
	//构造条件语句
	var whereArr []string

	if jobType != 0 {
		whereArr = append(whereArr, fmt.Sprintf("%s.type=%d", table, jobType))
	}
	if status != 4 {
		whereArr = append(whereArr, fmt.Sprintf("%s.status=%d", table, status))
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

	return where, nil
}

//GetList ...
func (oi *Job) GetList(paramsMap map[string]string) (*sql.Rows, error) {
	jobType := paramsMap["type"]
	jobStatus := paramsMap["status"]

	pageSize, err := strconv.Atoi(paramsMap["pageSize"])
	if err != nil {
		return nil, err
	}
	pageNum, err := strconv.Atoi(paramsMap["pageNum"])
	if err != nil {
		return nil, err
	}

	where, err := oi.buildConditions(jobType, jobStatus)
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

//GetCount ...
func (oi *Job) GetCount(paramsMap map[string]string) (int, error) {
	jobType := paramsMap["type"]
	jobStatus := paramsMap["status"]

	where, err := oi.buildConditions(jobType, jobStatus)
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
func (oi *Job) ScanRow(r *sql.Rows) (Record, error) {
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

//CreateJob ...
func CreateJob() (*Job, error) {
	obj := &Job{}
	err := obj.OpenDB("common.db.ssd_misc")
	if err != nil {
		return nil, err
	}
	err = obj.SetTableName("job")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
