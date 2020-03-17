package misc

import (
	// "fmt"
	"lib/serror"
	ssd_misc "model/misc"
)

//服务对象模型
type CanaryDevice struct{}

// GetDeviceIdList 获取灰度设备id列表
func (oi *CanaryDevice) GetDeviceIdList() ([]int, error) {
	//初始化灰度设备id列表
	deviceIDArr := make([]int, 0)

	//实例化数据模型
	canaryDeviceModel, err := ssd_misc.CreateCanaryDevice()
	// serror.Check(err)
	if err != nil {
		return deviceIDArr, err
	}

	defer func() {
		errClean := canaryDeviceModel.CloseDB()
		serror.Check(errClean)
	}()

	//查询灰度设备信息列表
	deviceList, err := canaryDeviceModel.GetList()
	// serror.Check(err)
	if err != nil {
		return deviceIDArr, err
	}

	defer deviceList.Close()

	// 遍历读取设备信息
	for deviceList.Next() {
		deviceInfo, err := canaryDeviceModel.ScanRow(deviceList)
		// serror.Check(err)
		if err != nil {
			return deviceIDArr, err
		}

		//计入当前灰度设备id列表
		deviceIDArr = append(deviceIDArr, deviceInfo.DeviceID)
	}

	return deviceIDArr, nil
}

//CreateCanaryDevice ...
func CreateCanaryDevice() (*CanaryDevice, error) {
	obj := &CanaryDevice{}
	return obj, nil
}
