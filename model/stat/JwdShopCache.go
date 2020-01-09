package stat

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

//JwdShopCacheRecord ...
//数据表字段结构
type JwdShopCacheRecord struct {
	ID       int    `json:"id"`
	ShopNo   string `json:"shop_no"`
	ShopName string `json:"shop_name"`
}

//JwdShopCache ... 数据模型对象
type JwdShopCache struct {
	model.Model
}

func (oi *JwdShopCache) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func (oi *JwdShopCache) prepare() *gorm.DB {
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

func (oi *JwdShopCache) getFields() []string {
	var keys []string
	var oRecord JwdShopCacheRecord
	dataType := reflect.TypeOf(oRecord)
	for k := 0; k < dataType.NumField(); k++ {
		keys = append(keys, dataType.Field(k).Tag.Get("json"))
	}
	return keys
}

//GetShopNoByName ...
// func (oi *JwdShopCache) GetShopNoByName(shopName string) (string, error) {
// 	shopNo := ""
// 	where := fmt.Sprintf("shop_name like '%s%s%s'", "%", shopName, "%")
// 	res := oi.prepare().Select("shop_no").Where(where).Row()
// 	err := res.Scan(&shopNo)
// 	if err != nil {
// 		return "", err
// 	}
// 	return shopNo, nil
// }

//GetListByName ...
func (oi *JwdShopCache) GetListByName(shopName string) (*sql.Rows, error) {
	var rows *sql.Rows
	where := fmt.Sprintf("shop_name like '%s%s%s' and shop_status=1", "%", shopName, "%")
	// limit := "5000"
	fields := oi.getFields()
	rows, err := oi.prepare().Select(fields).Where(where).Order("id asc").Rows()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

//GetListByShopNoArr ...
func (oi *JwdShopCache) GetListByShopNoArr(shopNoArr []string) (*sql.Rows, error) {
	var rows *sql.Rows

	if len(shopNoArr) == 0 {
		return rows, nil
	}

	shopNoStr := ""
	for k, v := range shopNoArr {
		if k == 0 {
			shopNoStr = fmt.Sprintf("'%s'", v)
		} else {
			shopNoStr = fmt.Sprintf("%s,'%s'", shopNoStr, v)
		}
	}

	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	//构造条件语句
	var whereArr []string

	where := ""

	whereArr = append(whereArr, fmt.Sprintf("shop_status=%d", 1))
	whereArr = append(whereArr, fmt.Sprintf("shop_no in (%s)", shopNoStr))

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

	// limit := "5000"
	fields := oi.getFields()

	// var rows *sql.Rows
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
func (oi *JwdShopCache) ScanRow(r *sql.Rows) (JwdShopCacheRecord, error) {
	var oRecord JwdShopCacheRecord
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

//CreateJwdShopCache ...
func CreateJwdShopCache() (*JwdShopCache, error) {
	obj := &JwdShopCache{}
	err := obj.OpenDB("common.db.ssd_stat")
	if err != nil {
		return nil, err
	}
	err = obj.SetTableName("jwd_shop_cache")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
