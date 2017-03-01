package main

import (
	"github.com/wupeaking/typhon/web"
	"github.com/wupeaking/typhon/app"
	"net/http"
)

type newhandler struct {
	web.BaseHandler
}

func (*newhandler) Get(w http.ResponseWriter, r *http.Request) error  {
	w.Write([]byte("这是一个新方法"))
	//var a map[string]string
	// 故意写一个异常
	//a["a"] = "a"
	return nil
}

func main()  {
	application := app.New(":9090", true)
	application.RegisterHandler("/", new(web.BaseHandler))
	application.RegisterHandler(`/get/\d+`, new(newhandler))
	application.Start()
}
