package tps

import (
	"database/sql"
	"fmt"
	// "github.com/huandu/xstrings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"lib/model"
	// "lib/stime"
	"crypto/md5"
	"github.com/huandu/xstrings"
	"io"
	"reflect"
	// "strings"
	// "strconv"
)

//数据表字段结构
type OrderInfoRecord struct {
	Id                   int    `json:"id"`
	Order_code           string `json:"order_code"`
	Pay_type             string `json:"pay_type"`
	Third_trade_no       string `json:"third_trade_no"`
	Pay_status           string `json:"pay_status"`
	Third_pay_status     string `json:"third_pay_status"`
	Reverse_status       string `json:"reverse_status"`
	Third_reverse_status string `json:"third_reverse_status"`
	Refund_status        string `json:"refund_status"`
	Third_refund_status  string `json:"third_refund_status"`
}

//数据模型对象
type OrderInfo struct {
	model.Model
	orderNo string
}

func (oi *OrderInfo) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

//设置订单号
func (oi *OrderInfo) SetOrderNo(orderNo string) error {
	clientId := "100001"
	oi.orderNo = clientId + orderNo
	fmt.Println(oi.orderNo)

	w := md5.New()
	io.WriteString(w, oi.orderNo)
	md5Str := fmt.Sprintf("%x", w.Sum(nil))
	fmt.Println(orderNo)
	fmt.Println(md5Str)
	md5SubStr := xstrings.Slice(md5Str, len(md5Str)-4, -1)
	fmt.Println(md5SubStr)

	asc2_1 := []byte(xstrings.Slice(md5SubStr, 0, 1))
	// if err != nil {
	// 	return err
	// }
	asc2_2 := []byte(xstrings.Slice(md5SubStr, 1, 2))
	// if err != nil {
	// 	return err
	// }
	asc2_3 := []byte(xstrings.Slice(md5SubStr, 2, 3))
	// if err != nil {
	// 	return err
	// }
	asc2_4 := []byte(xstrings.Slice(md5SubStr, 3, 4))
	// if err != nil {
	// 	return err
	// }
	fmt.Println(asc2_1)
	fmt.Println(asc2_2)
	hashNum := asc2_1[0] + asc2_2[0] + asc2_3[0] + asc2_4[0]
	fmt.Println(hashNum)
	tableNo := hashNum % 20
	table := fmt.Sprintf("order_%d", tableNo)
	fmt.Println(table)

	err := oi.SetTableName(table)
	if err != nil {
		return err
	}
	return nil
}

//获取当前要查询的数据所在的日期，格式：2019-10-12
func (oi *OrderInfo) GetOrderNo() (string, error) {
	return oi.orderNo, nil
}

func (oi *OrderInfo) prepare() *gorm.DB {
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

func (oi *OrderInfo) getFields() []string {
	var keys []string
	var orderInfoRecord OrderInfoRecord
	dataType := reflect.TypeOf(orderInfoRecord)
	for k := 0; k < dataType.NumField(); k++ {
		keys = append(keys, dataType.Field(k).Tag.Get("json"))
	}
	return keys
}

func (oi *OrderInfo) GetList(orderNo string) (*sql.Rows, error) {
	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	//构造条件语句
	var whereArr []string

	where := ""
	if orderNo != "" {
		whereArr = append(whereArr, fmt.Sprintf("order_code='%s'", orderNo))
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

	limit := "500"
	fields := oi.getFields()
	// fmt.Println(fields)
	// fmt.Println(where)

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

func (oi *OrderInfo) ScanRow(r *sql.Rows) (OrderInfoRecord, error) {
	var specInfoStruct OrderInfoRecord
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

func CreateOrderInfo() (*OrderInfo, error) {
	obj := &OrderInfo{}
	err := obj.OpenDB("common.db.ssd_tps")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
