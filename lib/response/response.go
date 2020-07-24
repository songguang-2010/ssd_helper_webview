package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ResponseError(w http.ResponseWriter, code int) {

	type resError struct {
		Msg string `json:"msg"`
	}

	res := &resError{
		Msg: "error occurred",
	}
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	body := string(b)

	//允许访问所有域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Cache-Control,Content-Language,Content-Type,Expires、Last-Modified,Pragma")

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//header的类型
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Connection, authorization, User-Agent, Cookie, token'")
	w.Header().Set("content-type", "application/json")
	// w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	//状态码
	w.WriteHeader(code)
	fmt.Fprint(w, body)
}

func ResponseSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Cache-Control,Content-Language,Content-Type,Expires、Last-Modified,Pragma")

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//header的类型
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Connection, authorization, User-Agent, Cookie, token'")

	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// w.Header().Set("X-Content-Type-Options", "nosniff")

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := string(b)

	// w.Header().Add("Last-Modified", "Thu, 18 Jun 2015 10:24:27 GMT")
	// w.Header().Add("Accept-Ranges", "bytes")
	// w.Header().Add("E-Tag", "55829c5b-17")
	// w.Header().Add("Server", "golang-http-server")
	// w.Write([]byte("<h1>\nHello world!\n</h1>\n"))
	w.Header().Set("Connection", "keep-alive")
	// w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	fmt.Fprint(w, body)
}
