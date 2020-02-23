module order

go 1.12

require (
	github.com/huandu/xstrings v1.2.0
	github.com/jinzhu/gorm v1.9.11
	github.com/spf13/viper v1.4.0
	lib/model v0.0.0
	lib/serror v0.0.0
	lib/stime v0.0.0
)

replace lib/serror v0.0.0 => ../../lib/serror

replace lib/model v0.0.0 => ../../lib/model

replace lib/stime v0.0.0 => ../../lib/stime
