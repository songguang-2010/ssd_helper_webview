module logwrap

go 1.12

require (
	github.com/spf13/viper v1.4.0
	lib/file v0.0.0
	lib/log v0.0.0
	lib/serror v0.0.0
)

replace lib/log v0.0.0 => ../../lib/log

replace lib/serror v0.0.0 => ../../lib/serror

replace lib/file v0.0.0 => ../../lib/file
