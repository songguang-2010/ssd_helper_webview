package aos

import (
	"database/sql"
	"fmt"
	// "github.com/huandu/xstrings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"lib/model"
	// "lib/stime"
	"github.com/huandu/xstrings"
	"reflect"
	"strings"
)

//数据表字段结构
type OrderInfoRecord struct {
	Id                int    `json:"id"`
	Order_no          string `json:"order_no"`
	Shop_id           string `json:"shop_id"`
	Shop_name         string `json:"shop_name"`
	Customer_phone    string `json:"customer_phone"`
	Goods_amount      string `json:"goods_amount"`
	Final_amount      string `json:"final_amount"`
	Activity_amount   string `json:"activity_amount"`
	Coupon_amount     string `json:"coupon_amount"`
	Discount_amount   string `json:"discount_amount"`
	Status            string `json:"status"`
	Pay_type          string `json:"pay_type"`
	Pay_amount        string `json:"pay_amount"`
	Pay_reduce_amount string `json:"pay_reduce_amount"`
	Refund_status     string `json:"refund_status"`
	Voucher_amount    string `json:"voucher_amount"`
}

//数据模型对象
type OrderInfo struct {
	model.Model
	currentDate string
}

func (oi *OrderInfo) getTableName() (string, error) {
	tableName, err := oi.Model.GetTableName()
	if err != nil {
		return "", err
	}
	return tableName, nil
}

//设置当前要查询的数据所在的日期，格式：2019-10-12
func (oi *OrderInfo) SetCurrentDate(currentDate string) error {
	oi.currentDate = currentDate
	dateNo := strings.Replace(currentDate, "-", "", -1)
	table := "order_info_" + xstrings.Slice(dateNo, 2, -1)

	err := oi.SetTableName(table)
	if err != nil {
		return err
	}
	return nil
}

//获取当前要查询的数据所在的日期，格式：2019-10-12
func (oi *OrderInfo) GetCurrentDate() (string, error) {
	return oi.currentDate, nil
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

func (oi *OrderInfo) GetList(shopName string, phone string) (*sql.Rows, error) {
	// len := len(shopId)
	// prefix := strings.Repeat(string('0'), (10 - len))
	//构造条件语句
	var whereArr []string

	where := ""
	if shopName != "" {
		whereArr = append(whereArr, fmt.Sprintf("shop_name like '%s%s%s'", "%", shopName, "%"))
	}
	if phone != "" {
		whereArr = append(whereArr, fmt.Sprintf("customer_phone=%s", phone))
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
	err := obj.OpenDB("common.db.ssd_aos")
	if err != nil {
		return nil, err
	}

	return obj, nil
}
