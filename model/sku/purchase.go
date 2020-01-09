package sku

import (
	"database/sql"
	"fmt"
	// "github.com/huandu/xstrings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"lib/model"
	// "lib/stime"
	"reflect"
	"strings"
)

//数据表字段结构
type PurchaseInfoRecord struct {
	ID            int    `json:"id"`
	Store_code    string `json:"store_code"`
	Store_name    string `json:"store_name"`
	Prod_code     string `json:"prod_code"`
	Prod_name     string `json:"prod_name"`
	Prod_number   string `json:"prod_number"`
	Date_request  string `json:"date_request"`
	Date_response string `json:"date_response"`
}

//数据模型对象
type PurchaseInfo struct {
	model.Model
	currentDate string
}

func (oi *PurchaseInfo) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func (oi *PurchaseInfo) prepare() *gorm.DB {
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

func (oi *PurchaseInfo) getFields() []string {
	var keys []string
	var oRecord PurchaseInfoRecord
	dataType := reflect.TypeOf(oRecord)
	for k := 0; k < dataType.NumField(); k++ {
		keys = append(keys, dataType.Field(k).Tag.Get("json"))
	}
	return keys
}

// ShopNoFilled ...
func (oi *PurchaseInfo) ShopNoFilled(shopNo string) string {
	strlen := len(shopNo)
	if strlen > 9 {
		return shopNo
	}
	prefix := strings.Repeat(string('0'), (10 - strlen))
	shopNoFilled := strings.TrimSpace(prefix + shopNo)
	return shopNoFilled
}

// ShopNoUnfilled ...
func (oi *PurchaseInfo) ShopNoUnfilled(shopNo string) string {
	shopNoUnfilled := strings.TrimLeft(shopNo, "0")
	return shopNoUnfilled
}

//GetRequestList ...
func (oi *PurchaseInfo) GetRequestList(shopNoArr []string, dateResponse string) (*sql.Rows, error) {
	//构造条件语句
	var whereArr []string

	where := ""

	whereArr = append(whereArr, fmt.Sprintf("type=0"))

	if len(shopNoArr) != 0 {
		shopNoStr := ""
		for k, v := range shopNoArr {
			shopNoFilled := oi.ShopNoFilled(v)
			if k == 0 {
				shopNoStr = fmt.Sprintf("'%s'", shopNoFilled)
			} else {
				shopNoStr = fmt.Sprintf("%s,'%s'", shopNoStr, shopNoFilled)
			}
		}
		// fmt.Println(shopNoStr)
		whereArr = append(whereArr, fmt.Sprintf("store_code in (%s)", shopNoStr))
	}

	if dateResponse != "" {
		whereArr = append(whereArr, fmt.Sprintf("date_response='%s'", strings.TrimSpace(dateResponse)))
	}

	if len(whereArr) > 0 {
		for k, v := range whereArr {
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
		rows, err = oi.prepare().Select(fields).Where(where).Order("id asc").Limit(limit).Rows()
	} else {
		rows, err = oi.prepare().Select(fields).Order("id asc").Limit(limit).Rows()
	}

	if err != nil {
		return nil, err
	}
	return rows, nil
}

//GetResponseList ...
func (oi *PurchaseInfo) GetResponseList(shopNoArr []string, dateResponse string) (*sql.Rows, error) {
	//构造条件语句
	var whereArr []string

	where := ""

	whereArr = append(whereArr, fmt.Sprintf("type=1"))

	if len(shopNoArr) != 0 {
		shopNoStr := ""
		for k, v := range shopNoArr {
			len := len(v)
			prefix := strings.Repeat(string('0'), (10 - len))
			shopNoFilled := strings.TrimSpace(prefix + v)
			if k == 0 {
				shopNoStr = fmt.Sprintf("'%s'", shopNoFilled)
			} else {
				shopNoStr = fmt.Sprintf("%s,'%s'", shopNoStr, shopNoFilled)
			}
		}
		fmt.Println(shopNoStr)
		whereArr = append(whereArr, fmt.Sprintf("store_code in (%s)", shopNoStr))
	}

	if dateResponse != "" {
		whereArr = append(whereArr, fmt.Sprintf("date_response='%s'", strings.TrimSpace(dateResponse)))
	}

	if len(whereArr) > 0 {
		for k, v := range whereArr {
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
		rows, err = oi.prepare().Select(fields).Where(where).Order("id asc").Limit(limit).Rows()
	} else {
		rows, err = oi.prepare().Select(fields).Order("id asc").Limit(limit).Rows()
	}

	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (oi *PurchaseInfo) ScanRow(r *sql.Rows) (PurchaseInfoRecord, error) {
	var recordStruct PurchaseInfoRecord
	db, err := oi.GetDB()
	if err != nil {
		oi.SetError(err.Error())
		return recordStruct, nil
	}
	err = db.ScanRows(r, &recordStruct)
	if err != nil {
		return recordStruct, err
	}
	return recordStruct, nil
}

//CreatePurchaseInfo ...
func CreatePurchaseInfo() (*PurchaseInfo, error) {
	obj := &PurchaseInfo{}
	err := obj.OpenDB("common.db.ssd_sku")
	if err != nil {
		return nil, err
	}
	err = obj.SetTableName("shop_goods_purchase")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
