## 一个使用golang编写的简易的web框架 灵感来自于Python的tornado

### 一个简单的示例

```shell
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
	h.Write([]byte(`this is post response`))
	fmt.Println(h.GetArguments())
	fmt.Println(h.GroupArgs)
	return nil
}


func main()  {
	application := app.New(":9090", true)
	application.RegisterRouter("/", new(web.BaseHandler))
	# 注册正则路由
	application.RegisterRouter(`/get/\d+`, new(newhandler))
        # 注册分组的正则路由
        application.RegisterRouter(`/aaa/(?P<name>\w+)/(?P<id>\d+)`, new(newhandler))
	application.Start()
}
```
