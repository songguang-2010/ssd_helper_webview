package stime

import (
	"fmt"
	"time"
)

const (
	//时间转换的模板，golang里面只能是 "2006-01-02 15:04:05" （go的诞生时间）
	tpl_datetime = "2006-01-02 15:04:05"
)

func GetTsByStr(str string) (int64, error) {
	//使用parseInLocation将字符串格式化返回本地时区时间对象
	timeObj, _ := time.ParseInLocation(tpl_datetime, str, time.Local)
	return timeObj.Unix(), nil
}

func GetTsRangeByHour(date string, h int) (int64, int64, error) {
	datetimeStart := date + fmt.Sprintf(" %.2d:00:00", h)
	datetimeEnd := date + fmt.Sprintf(" %.2d:59:59", h)
	timestampStart, _ := GetTsByStr(datetimeStart)
	timestampEnd, _ := GetTsByStr(datetimeEnd)
	return timestampStart, timestampEnd, nil
}
