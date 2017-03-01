package app

import (
	"net/http"
	"regexp"
	"github.com/wupeaking/typhon/app/utils"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"runtime"
)

// 定义路由接口
type RouteHandler interface {
	Get(w http.ResponseWriter, r *http.Request) error
	Post(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
	Put(w http.ResponseWriter, r *http.Request) error
	Patch(w http.ResponseWriter, r *http.Request) error
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
func (app *Application)RegisterHandler(url string, handler RouteHandler) {
	app.entry[url] = handler
	// 为每个注册的entry生成一个正则 加速URL匹配
	app.entryRegex[url], _ = regexp.Compile(url)
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
				error: %s
				`, r.URL.Path, r.Method, body, err)
		}
		if err != nil && !app.debug {
			fmt.Fprintf(w, "Internal service error")
		}
	}()
	endpoint := r.URL.Path
	isMatch := false
	var handler RouteHandler
	for url, h := range app.entry{
		urlRe := app.entryRegex[url]
		if !urlRe.MatchString(endpoint){
			continue
		}else{
			isMatch = true
			handler = h
		}
	}
	if !isMatch{
		log.WithField("url", endpoint).Debug("no method can match the url")
		w.Write(utils.Str2bytes("error endpoint"))
		return
	}
	switch r.Method {
	case "GET":
		handler.Get(w, r)
	case "POST":
		handler.Post(w, r)
	case "DELET":
		handler.Delete(w, r)
	case "PATCH":
		handler.Patch(w, r)
	}
}

func(app *Application)Start() error{
	e := http.ListenAndServe(":9090", app)
	if e != nil {
		log.WithField("err:", e).Fatal("start web server fail")
	}
	return e
}