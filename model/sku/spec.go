package sku

import (
	"database/sql"
	"fmt"
	// "github.com/huandu/xstrings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"lib/model"
	// "lib/stime"
	// "reflect"
	"strings"
)

//数据表字段结构
type SpecInfoRecord struct {
	Id                  int    `json:"id"`
	Store_code          string `json:"store_code"`
	Store_name          string `json:"store_name"`
	Prod_code           string `json:"prod_code"`
	Prod_name           string `json:"prod_name"`
	Sale_unit           string `json:"sale_unit"`
	Purchase_unit       string `json:"purchase_unit"`
	Sale_unit_ratio     string `json:"sale_unit_ratio"`
	Purchase_unit_ratio string `json:"purchase_unit_ratio"`
}

//数据模型对象
type SpecInfo struct {
	model.Model
	currentDate string
}

func (oi *SpecInfo) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

//设置当前要查询的数据所在的日期，格式：2019-10-12
func (oi *SpecInfo) SetCurrentDate(currentDate string) error {
	oi.currentDate = currentDate
	dateNo := strings.Replace(currentDate, "-", "", -1)
	table := "shop_goods_spec_" + dateNo

	err := oi.SetTableName(table)
	if err != nil {
		return err
	}
	return nil
}

//获取当前要查询的数据所在的日期，格式：2019-10-12
func (oi *SpecInfo) GetCurrentDate() (string, error) {
	return oi.currentDate, nil
}

func (oi *SpecInfo) prepare() *gorm.DB {
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

// GetList ...
func (oi *SpecInfo) GetList(shopCode string, shopName string, prodName string, limit int) (*sql.Rows, error) {
	shopCodeLen := len(shopCode)
	prefix := strings.Repeat(string('0'), (10 - shopCodeLen))
	//构造条件语句
	where := ""
	var whereArr []string
	if shopName != "" {
		whereArr = append(whereArr, fmt.Sprintf("store_name like '%s%s%s'", "%", shopName, "%"))
	}
	if prodName != "" {
		whereArr = append(whereArr, fmt.Sprintf("prod_name like '%s%s%s'", "%", prodName, "%"))
	}
	if shopCode != "" {
		whereArr = append(whereArr, fmt.Sprintf("store_code='%s%s'", prefix, shopCode))
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

	fields := "id, store_code, store_name, prod_code, prod_name, sale_unit, purchase_unit, sale_unit_ratio, purchase_unit_ratio"

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

func (oi *SpecInfo) ScanRow(r *sql.Rows) (SpecInfoRecord, error) {
	var specInfoStruct SpecInfoRecord
	db, err := oi.GetDB()
	if err != nil {
		oi.SetError(err.Error())
		return specInfoStruct, nil
	}
	err = db.ScanRows(r, &specInfoStruct)
	if err != nil {
		return specInfoStruct, err
	}
	return specInfoStruct, nil
}

func CreateSpecInfo() (*SpecInfo, error) {
	specInfo := &SpecInfo{}
	err := specInfo.OpenDB("common.db.ssd_sku")
	if err != nil {
		return nil, err
	}

	return specInfo, nil
}
