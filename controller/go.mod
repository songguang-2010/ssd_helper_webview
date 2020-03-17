module controller

go 1.12

require (
	lib/model v0.0.0
	lib/response v0.0.0
	lib/serror v0.0.0
	lib/stime v0.0.0
	model/aos v0.0.0
	model/misc v0.0.0
	model/order v0.0.0
	model/sku v0.0.0
	model/stat v0.0.0
	model/tps v0.0.0
	service/misc v0.0.0
)

replace lib/response v0.0.0 => ../lib/response

replace lib/serror v0.0.0 => ../lib/serror

replace lib/model v0.0.0 => ../lib/model

replace lib/stime v0.0.0 => ../lib/stime

replace model/sku v0.0.0 => ../model/sku

replace model/order v0.0.0 => ../model/order

replace model/aos v0.0.0 => ../model/aos

replace model/tps v0.0.0 => ../model/tps

replace model/misc v0.0.0 => ../model/misc

replace model/stat v0.0.0 => ../model/stat

replace service/misc v0.0.0 => ../service/misc
