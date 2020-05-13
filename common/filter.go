package common

import "net/http"

// 声明一个函数类型
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type Filter struct {
	// 用于存储对应 uri  的 FilterHanler 函数
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

// 注册拦截器
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// 获取拦截器
func (f *Filter) GetFilterUri(uri string, handler FilterHandle) FilterHandle {
	return f.filterMap[uri]
}

type WebHandler func(rw http.ResponseWriter, req *http.Request)

// 注册拦截器
func (f *Filter) Handler(webHandler WebHandler) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		for path, handle := range f.filterMap {
			// 如果匹配则只需拦截业务逻辑
			if path == req.RequestURI {
				if err := handle(rw, req); err != nil {
					rw.Write([]byte(err.Error()))
					// 执行错误，直接退出不执行后续内容
					return
				}
				// 匹配执行拦截完成，跳出循环
				break
			}
		}

		// 执行正常业务逻辑
		webHandler(rw, req)
	}
}
