module coroutine

go 1.12

require (
	lib/file v0.0.0
	lib/log v0.0.0
	lib/logwrap v0.0.0
	lib/serror v0.0.0
)

replace lib/logwrap v0.0.0 => ../../lib/logwrap

replace lib/serror v0.0.0 => ../../lib/serror

replace lib/file v0.0.0 => ../../lib/file

replace lib/log v0.0.0 => ../../lib/log
