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
	"github.com/wupeaking/typhon/web"
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
	// 注册默认的根路径handler
	app.RegisterRouter("/", new(routerControl))
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

// 编写根路径的默认handler

// 默认的首页html
var html = `<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <title>typhone index</title>
	<!-- 最新版本的 Bootstrap 核心 CSS 文件 -->
	<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

	<!-- 可选的 Bootstrap 主题文件（一般不用引入） -->
	<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

	<!-- 最新的 Bootstrap 核心 JavaScript 文件 -->
	<script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
  </head>
  <body>
    <h1>你好，世界！</h1>
    <h3>这是一个简单的web框架 基于golang编写 用于自己的小项目中</h3>
	<h6>下面是一个简单的示例: </h6>
    <div>
    	<pre>
    		package main
			import (
				"github.com/wupeaking/typhon/web"
				"github.com/wupeaking/typhon/app"
				"net/http"
			)

			type newhandler struct {
				web.BaseHandler
			}

			// 打印异常信息栈 获取请求的参数 以及正则匹配的分组参数
			func (h *newhandler) Get() error  {
				fmt.Println("分组内容: ", h.GroupArgs)
				h.Write([]byte("这是一个新方法"))

				println(h.GetArguments())
				var a map[string]string
				// 故意写一个异常
				a["a"] = "a"
				return nil
			}

			// 获取post方法的请求参数 分组参数
			func (h *newhandler) Post() error {
				h.Write([]byte("this is post response"))
				fmt.Println(h.GetArguments())
				fmt.Println(h.GroupArgs)
				return nil
			}

			func main()  {
				application := app.New(":9090", true)
				application.RegisterRouter("/", new(web.BaseHandler))
				# 注册正则路由
				application.RegisterRouter("/get/\d+", new(newhandler))
					# 注册分组的正则路由
					application.RegisterRouter("/aaa/(?P<name>\w+)/(?P<id>\d+)", new(newhandler))
				application.Start()
			}
    	</pre>
    </div>

    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="https://cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
  </body>
</html>`

type routerControl struct {
	web.BaseHandler
}

func (ctl *routerControl)Get() error {
	ctl.Write(utils.Str2bytes(html))
	return nil
}