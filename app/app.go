package app

import (
	"net/http"
	"regexp"
	"github.com/wupeaking/typhon/app/utils"
	//"github.com/wupeaking/typhon/web"
	log "github.com/Sirupsen/logrus"
	deepcopy "github.com/mitchellh/copystructure"
	"fmt"
	"runtime"
	"reflect"
)

// 定义路由接口
type RouteHandler interface {
	Get() error
	Post() error
	Delete() error
	Put() error
	Patch() error
	Finish() error
}

type Application struct {
	// 自定义自己的路由
	http.ServeMux
	// 监听地址
	addr string
	// 是否是debug模式
	debug bool
	entry map[string]RouteHandler
	entryRegex map[string]*regexp.Regexp

}

func New(addr string, debug bool) *Application{
	app := &Application{addr:addr, debug:debug}
	if debug{
		log.SetLevel(log.DebugLevel)
	}
	app.entry = make(map[string]RouteHandler)
	app.entryRegex = make(map[string]*regexp.Regexp)
	return app
}

// 注册全局endpoint<--->handler
func (app *Application)RegisterRouter(url string, handler RouteHandler) {
	url = "^" + url + "$"
	app.entry[url] = handler
	// 为每个注册的entry生成一个正则 加速URL匹配
	app.entryRegex[url] = regexp.MustCompile(url)
}

// 实现http.Handler的接口定义方法
func (app *Application)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 捕获错误异常
	defer func (){
		err := recover()
		if err != nil{
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			//_, file, line, _  := runtime.Caller(0)
			log.WithFields(log.Fields{
				"error": err,
			}).Error("catch handler exception", utils.Bytes2str(buf[:n]))
		}
		if err != nil && app.debug {
			var body string
			buf := make([]byte, 1024)
			len, _ := r.Body.Read(buf)
			if len >= 1024{
				len = 1024
			}else if len == 0 {
				len = 1
			}
			body = utils.Bytes2str(buf[0: len-1])
			fmt.Fprintf(w, `catch hander exception:
	url: %s
	method: %s
	request body: %s
	error: %s`, r.URL.Path, r.Method, body, err)
		}
		if err != nil && !app.debug {
			fmt.Fprintf(w, "Internal service error")
		}
	}()

	endpoint := r.URL.Path
	isMatch := false
	var handler RouteHandler
	// url正则匹配后传递给handler的分组参数
	var reMatchArgs map[string]string
	for url, h := range app.entry{
		urlRe := app.entryRegex[url]
		if !urlRe.MatchString(endpoint){
			continue
		}else{
			isMatch = true
			handler = h
			// 解析出分组 传递给handler 获取所有的分组名["", "分组名1", "分组名2"]
			groupNames := urlRe.SubexpNames()
			// 获取所有的分组个数
			groupNum := urlRe.NumSubexp()

			// 判断分组名的个数是否和分组数量一致 如果不一致 提出警告 不然匹配的分组内容和分组名称不一致
			if len(groupNames)-1 != groupNum {
				log.Warn("url re match group name count not eq group count")
			}
			if len(groupNames) > 1 {
				reMatchArgs = make(map[string]string)
				// 查找所有的分组
				allMatch := urlRe.FindStringSubmatch(endpoint)
				i := 0
				for _, groupName := range groupNames{
					if i == 0{
						i++
						continue
					}
					reMatchArgs[groupName] = allMatch[i]
					i++
				}
			}

		}
	}
	if !isMatch{
		log.WithField("url", endpoint).Debug("no method can match the url")
		w.Write(utils.Str2bytes("error endpoint"))
		return
	}

	// 创建一个新的control handler
	newobji, _ := deepcopy.Copy(handler)
	// 强制转换为web.BaseHandler
	//if reflect.TypeOf(newobji).Kind() == reflect.Ptr{
	//	newobj := newobji.(*web.BaseHandler)
	//	newobj.Request = r
	//	newobj.Response = w
	//	newobj.GroupArgs = reMatchArgs
	//}else{
	//	newobj := newobji.(web.BaseHandler)
	//	newobj.Request = r
	//	newobj.Response = w
	//	newobj.GroupArgs = reMatchArgs
	//}
	// 为新的对象赋值
	reflect.ValueOf(newobji).Elem().FieldByName("Response").Set(reflect.ValueOf(w))
	reflect.ValueOf(newobji).Elem().FieldByName("Request").Set(reflect.ValueOf(r))
	reflect.ValueOf(newobji).Elem().FieldByName("GroupArgs").Set(reflect.ValueOf(reMatchArgs))
	handler = newobji.(RouteHandler)

	switch r.Method {
	case "GET":
		handler.Get()
	case "POST":
		handler.Post()
	case "DELET":
		handler.Delete()
	case "PATCH":
		handler.Patch()
	}
	handler.Finish()
}

func(app *Application)Start() error{
	e := http.ListenAndServe(app.addr, app)
	if e != nil {
		log.WithField("err:", e).Fatal("start web server fail")
	}
	return e
}