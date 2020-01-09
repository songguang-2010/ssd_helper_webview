package serror

// custom error object
type CustomError struct {
	msg string
}

// override the error method for interface error
func (ce *CustomError) Error() string {
	return ce.msg
}

// check function for global
func Check(err error) {
    if err != nil {
        panic(err)
    }
}

// init custom error object
func New(msg string) *CustomError {
	return &CustomError{msg}
}
