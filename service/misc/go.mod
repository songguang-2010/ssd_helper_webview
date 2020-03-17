module misc

go 1.12

require (
	lib/model v0.0.0
	lib/serror v0.0.0
	model/misc v0.0.0
)

replace lib/serror v0.0.0 => ../lib/serror

replace lib/model v0.0.0 => ../lib/model

replace model/misc v0.0.0 => ../model/misc
