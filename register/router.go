package register

import (
	"lib/route"
)

// InitRouter ...
func InitRouter() {
	//添加路由项
	route.AddRoute(route.RouteItem{
		Method:     "POST",
		Path:       "/set-canary",
		Controller: "WebController",
		Function:   "SetCanary",
	})
	route.AddRoute(route.RouteItem{
		Method:     "POST",
		Path:       "/login",
		Controller: "WebController",
		Function:   "Login",
	})
	route.AddRoute(route.RouteItem{
		Method:     "GET",
		Path:       "/get-sku-specs",
		Controller: "WebController",
		Function:   "GetSkuSpecs",
	})
	route.AddRoute(route.RouteItem{
		Method:     "GET",
		Path:       "/get-sku-requests",
		Controller: "WebController",
		Function:   "GetSkuRequests",
	})
	route.AddRoute(route.RouteItem{
		Method:     "GET",
		Path:       "/get-sku-responses",
		Controller: "WebController",
		Function:   "GetSkuResponses",
	})
	route.AddRoute(route.RouteItem{
		Method:     "GET",
		Path:       "/get-ssd-orders",
		Controller: "WebController",
		Function:   "GetSsdOrders",
	})
	route.AddRoute(route.RouteItem{
		Method:     "GET",
		Path:       "/get-aos-orders",
		Controller: "WebController",
		Function:   "GetAosOrders",
	})
	route.AddRoute(route.RouteItem{
		Method:     "GET",
		Path:       "/get-tps-orders",
		Controller: "WebController",
		Function:   "GetTpsOrders",
	})
	route.AddRoute(route.RouteItem{
		Method:     "GET",
		Path:       "/get-misc-devices",
		Controller: "WebController",
		Function:   "GetMiscDevices",
	})
}
