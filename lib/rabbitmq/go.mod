module rabbitmq

go 1.12

require (
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	lib/coroutine v0.0.0
	lib/file v0.0.0
	lib/log v0.0.0
	lib/logwrap v0.0.0
	lib/serror v0.0.0
)

replace lib/log v0.0.0 => ../../lib/log

replace lib/logwrap v0.0.0 => ../../lib/logwrap

replace lib/file v0.0.0 => ../../lib/file

replace lib/serror v0.0.0 => ../../lib/serror

replace lib/coroutine v0.0.0 => ../../lib/coroutine
