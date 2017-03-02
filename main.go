package main

import (
	"github.com/wupeaking/typhon/web"
	"github.com/wupeaking/typhon/app"
	"fmt"
)

type newhandler struct {
	web.BaseHandler
}

func (h *newhandler) Get() error  {
	fmt.Println("分组内容: ", h.GroupArgs)
	h.Write([]byte("这是一个新方法"))

	println(h.GetArguments())
	var a map[string]string
	// 故意写一个异常
	a["a"] = "a"
	return nil
}

func (h *newhandler) Post() error {
	h.Write([]byte(`this is post response`))
	fmt.Println(h.GetArguments())
	fmt.Println(h.GroupArgs)
	return nil
}

func main()  {
	application := app.New(":9090", true)
	application.RegisterRouter("/", new(web.BaseHandler))
	application.RegisterRouter(`/get/(?P<name>\d+)`, new(newhandler))
	application.RegisterRouter(`/aaa/(?P<name>\w+)/(?P<id>\d+)`, new(newhandler))
	application.Start()
}
