package register

import (
	"lib/route"
	"middleware"
)

// InitMiddleware ...
func InitMiddleware() {
	// 注册中间件
	route.AddMiddleware("WebMiddleware", &middleware.WebMiddleware{})
}
