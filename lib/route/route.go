package route

import (
	"fmt"
	"lib/file"
	// "path/filepath"
	// "lib/log"
	"net/http"
	// "os"
	"reflect"
)

// type route struct {
// 	path       string   // 路径
// 	target     string   // 对应的控制器路径 Controller@index 这样的方法
// 	method     string   // 访问类型 是get post 或者其他
// 	alias      string   // 路由的别名
// 	middleware []string // 中间件名称
// 	controller string   // 控制器名称
// 	function   string   // 挂载到控制器上的方法名称
// }

// type route_group struct {
// 	root_path   string   // 路径
// 	root_target string   // 对应的控制器路径 Controller@index 这样的方法
// 	alias       string   // 路由的别名
// 	middleware  []string // 中间件名称
// 	routes      []route  // 包含的路由
// }

type RouteItem struct {
	Path       string   // 路径
	Method     string   // 访问类型 是get post 或者其他
	Middleware []string // 中间件名称
	Controller string   // 控制器名称
	Function   string   // 挂载到控制器上的方法名称
}

//服务监听控制器结构
type MyMux struct {
	// 静态文件目录
	publicPath string
}

// 单个的路由集合
// var Routes []route
var Routes []RouteItem

// 路由组集合
// var RoutesGroups []route_group

// 全部路由列表
// var RoutesList []route
var RoutesList []RouteItem

// var R interface{}

// 这里记录所有的应该注册的结构体

// 控制器map
var ControllerMap map[string]interface{}

// 中间件map
var MiddlewareMap map[string]interface{}

// init ...
func init() {
	ControllerMap = make(map[string]interface{})
	MiddlewareMap = make(map[string]interface{})
}

// 设置静态目录
func (p *MyMux) SetPublicPath(publicPath string) {
	p.publicPath = string(publicPath)
}

// 获取静态目录
func (p *MyMux) GetPublicPath() string {
	return p.publicPath
}

//服务监听控制器方法
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("current path : ", r.URL.Path)
	fmt.Println("current method : ", r.Method)
	//返回数据格式是json
	r.ParseForm()
	// fmt.Println("收到客户端请求: ", r.Form)
	// fmt.Println("收到客户端请求: ", r.Form.Get("title"))
	// 这个变量用来标记是否找到了动态路由
	flag := false
	// 每一个http请求都会走到这里，然后在这里，根据请求的URL，为其分配所需要调用的方法
	params := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
	for _, v := range RoutesList {
		// 检测路由，根据路由指向需要的数据
		if r.URL.Path == v.Path && r.Method == v.Method {
			// 寻找到了对应路由，无需使用静态服务器
			flag = true

			// 检测该路由中是否存在中间件，如果存在，顺序调用
			for _, m := range v.Middleware {
				// 判断是否注册了这个中间件
				if mid, ok := MiddlewareMap[m]; ok {
					rmid := reflect.ValueOf(mid)
					// 执行中间件，返回values数组
					params = rmid.MethodByName("Handle").Call(params)
					// 判断中间件执行结果，是否还要继续往下走
					str := rmid.Elem().FieldByName("ResString").String()
					if str != "" {
						status := rmid.Elem().FieldByName("Status").Int()
						// 字符串不空，查看状态码，默认返回500错误
						if status == 0 {
							status = 500
						}
						w.WriteHeader(int(status))
						fmt.Fprint(w, str)

						return
					}
				}
			}

			// 检测成功，开始调用方法
			// 获取一个控制器包下的结构体
			// 存在  c为结构体，调用c上挂载的方法
			if d, ok := ControllerMap[v.Controller]; ok {
				reflect.ValueOf(d).MethodByName(v.Function).Call(params)
			}

			// 停止向后执行
			return
		}
	}

	// 如果路由列表中还是没有的话,去静态服务器中寻找
	if !flag {
		// 去静态目录中寻找
		publicPath := p.GetPublicPath() + r.URL.Path
		ok, _ := file.PathExists(publicPath)
		if ok {
			http.ServeFile(w, r, publicPath)
		} else {
			http.NotFound(w, r)
		}

	}

	return
}

//AddRoute 添加单个路由
func AddRoute(route RouteItem) {
	// 把这个新的路由，放到单个路由切片中，也要放到路由列表中
	Routes = append(Routes, route)
	RoutesList = append(RoutesList, route)
}

//AddMiddleWare 添加单个中间件
func AddMiddleware(midName string, midObj interface{}) {
	MiddlewareMap[midName] = midObj
}

//AddController 添加单个控制器
func AddController(controllerName string, controllerObj interface{}) {
	ControllerMap[controllerName] = controllerObj
}

// func init() {
// 	// 初始化方法，加载路由文件
// 	// 获取路由路径，根据路由路径获取所有路由文件，然后读取所有文件，赋值给当前成员变量
// 	routes_path := GetRoutesPath()
// 	dir_list, err := ioutil.ReadDir(routes_path)
// 	Check(err)
// 	// 根据dir list 遍历所有文件 获取所有json文件，拿到所有的路由 路由组
// 	for _, v := range dir_list {
// 		fmt.Println("正在加载路由文件........" + v.Name())
// 		// 读取文件内容，转换成json，并且加入数组中
// 		content, err := FileGetContents(routes_path + "/" + v.Name())
// 		Check(err)
// 		err = json.Unmarshal([]byte(content), &R)
// 		Check(err)
// 		// 开始解析R,将其分类放入全局变量中
// 		parse(R)
// 	}
// }

// func parse(r interface{}) {
// 	// 拿到了r 我们要解析成实际的数据
// 	m := r.(map[string]interface{})
// 	//newRoute := route{}
// 	for k, v := range m {
// 		if k == "Routes" {
// 			// 解析单个路由
// 			parseRoutes(v)
// 		}
// 		if k == "RoutesGroups" {
// 			// 解析路由组
// 			parseRoutesGroups(v)
// 		}
// 	}

// }

// // 添加一个单一路由

// // 解析json文件中的单一路由的集合
// func parseRoutes(r interface{}) {
// 	m := r.([]interface{})
// 	for _, v := range m {
// 		// v 就是单个的路由了
// 		simpleRoute := v.(map[string]interface{})
// 		// 定义一个路由结构体
// 		newRoute := route{}
// 		for kk, vv := range simpleRoute {
// 			switch kk {
// 			case "Route":
// 				newRoute.path = vv.(string)
// 				break
// 			case "Target":
// 				newRoute.target = vv.(string)
// 				break
// 			case "Method":
// 				newRoute.method = vv.(string)
// 				break
// 			case "Alias":
// 				newRoute.alias = vv.(string)
// 				break
// 			case "Middleware":
// 				//newRoute.middleware = vv.([])
// 				var mdw []string
// 				vvm := vv.([]interface{})
// 				for _, vvv := range vvm {
// 					mdw = append(mdw, vvv.(string))
// 				}
// 				newRoute.middleware = mdw
// 				break
// 			default:
// 				break
// 			}
// 		}

// 		// 把target拆分成控制器和方法
// 		cf := strings.Split(newRoute.target, "@")
// 		if len(cf) == 2 {
// 			newRoute.controller = cf[0]
// 			newRoute.function = cf[1]
// 		} else {
// 			fmt.Println("Target格式错误！" + newRoute.target)
// 			return
// 		}

// 		// 把这个新的路由，放到单个路由切片中，也要放到路由列表中

// 		Routes = append(Routes, newRoute)
// 		RoutesList = append(RoutesList, newRoute)
// 	}
// }

// func parseRoutesGroups(r interface{}) {
// 	// 解析路由组
// 	m := r.([]interface{})
// 	for _, v := range m {
// 		group := v.(map[string]interface{})
// 		for kk, vv := range group {
// 			// 新建一个路由组结构体
// 			var newGroup route_group
// 			switch kk {
// 			case "RootRoute":
// 				newGroup.root_path = vv.(string)
// 				break
// 			case "RootTarget":
// 				newGroup.root_target = vv.(string)
// 				break
// 			case "Middleware":
// 				var mdw []string
// 				vvm := vv.([]interface{})
// 				for _, vvv := range vvm {
// 					mdw = append(mdw, vvv.(string))
// 				}
// 				newGroup.middleware = mdw
// 				break
// 			case "Routes":
// 				// 由于涉及到根路由之类的概念，所以不能使用上面的parseRoutes方法，需要再写一个方法用来解析真实路由
// 				rs := parseRootRoute(group)
// 				newGroup.routes = rs
// 				break
// 			default:
// 				break
// 			}
// 			// 把这个group放到路由组里
// 			RoutesGroups = append(RoutesGroups, newGroup)
// 		}
// 	}
// }

// // 解析根路由 传入根路由路径 目标跟路径 并且传入路由inteface列表，返回一个完整的路由集合
// // 只传入一个路由组，返回一个完整的路由集合
// func parseRootRoute(group map[string]interface{}) []route {
// 	// 获取路由根路径和目标根路径,还有公共中间件
// 	var tmpRoutes []route // 要返回的路由切片
// 	var route_root_path string
// 	var target_root_path string
// 	var public_middleware []string
// 	for k, v := range group {
// 		if k == "RootRoute" {
// 			route_root_path = v.(string)
// 		}
// 		if k == "RootTarget" {
// 			target_root_path = v.(string)
// 		}
// 		if k == "Middleware" {
// 			vvm := v.([]interface{})
// 			for _, vvv := range vvm {
// 				public_middleware = append(public_middleware, vvv.(string))
// 			}
// 		}
// 	}

// 	// 开始获取路由
// 	for k, s := range group {
// 		if k == "Routes" {
// 			m := s.([]interface{})
// 			for _, v := range m {
// 				// v 就是单个的路由了
// 				simpleRoute := v.(map[string]interface{})
// 				// 定义一个路由结构体
// 				newRoute := route{}
// 				for kk, vv := range simpleRoute {
// 					switch kk {
// 					case "Route":
// 						newRoute.path = route_root_path + vv.(string)
// 						break
// 					case "Target":
// 						newRoute.target = target_root_path + vv.(string)
// 						break
// 					case "Method":
// 						newRoute.method = vv.(string)
// 						break
// 					case "Alias":
// 						newRoute.alias = vv.(string)
// 						break
// 					case "Middleware":
// 						vvm := vv.([]interface{})
// 						for _, vvv := range vvm {
// 							newRoute.middleware = append(public_middleware, vvv.(string)) // 公共的和新加入的放在一起就是总共的
// 						}

// 						break
// 					default:
// 						break
// 					}
// 				}
// 				// 把target拆分成控制器和方法
// 				cf := strings.Split(newRoute.target, "@")
// 				if len(cf) == 2 {
// 					newRoute.controller = cf[0]
// 					newRoute.function = cf[1]
// 				} else {
// 					fmt.Println("Target格式错误！" + newRoute.target)
// 					os.Exit(2)
// 				}
// 				// 把这个新的路由，放到路由列表中，并且返回放到路由集合中，作为返回值返回
// 				RoutesList = append(RoutesList, newRoute)
// 				tmpRoutes = append(tmpRoutes, newRoute)
// 			}
// 		}
// 	}
// 	return tmpRoutes
// }
