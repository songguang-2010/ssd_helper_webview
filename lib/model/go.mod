module model

go 1.12

require (
	github.com/jinzhu/gorm v1.9.11
	github.com/spf13/viper v1.4.0
	lib/serror v0.0.0
)

replace lib/serror v0.0.0 => ../../lib/serror
