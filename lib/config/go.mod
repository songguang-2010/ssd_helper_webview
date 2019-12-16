module config

go 1.12

require (
	github.com/spf13/viper v1.4.0
	golang.org/x/text v0.3.2 // indirect
	lib/file v0.0.0
	lib/serror v0.0.0
)

replace lib/file v0.0.0 => ../../lib/file

replace lib/serror v0.0.0 => ../../lib/serror
