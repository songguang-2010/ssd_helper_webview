moudle register

go 1.12

require (
	controller v0.0.0
	lib/route v0.0.0
	middleware v0.0.0
)

replace controller v0.0.0 => ../controller

replace middleware v0.0.0 => ../middleware

replace lib/route v0.0.0 => ../lib/route