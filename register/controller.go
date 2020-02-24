package register

import (
	"controller"
	"lib/route"
)

// InitController ...
func InitController() {
	// 注册控制器
	// ControllerMap["WebController"] = &controller.WebController{}
	route.AddController("WebController", &controller.WebController{})
}
