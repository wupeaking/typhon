/*
* 时间: 2017年03月01日13:12:19
* 说明: 一个简单的web框架 仿照tornado的框架模型
* 作者: wupeaking
*/

package web

import (
	"net/http"
	"github.com/wupeaking/typhon/app/utils"
	json "github.com/bitly/go-simplejson"
	log "github.com/Sirupsen/logrus"
	"time"
	"bytes"
	"strings"
	//"fmt"
)

// 一个基本的handler应该至少实现http的PUT DELET POST GET PATCH五种方法
type BaseHandler struct {
	Response http.ResponseWriter
	Request *http.Request
	GroupArgs map[string]string
	Arguments map[string]interface{}
	buf bytes.Buffer
	isfinish bool
}

func (handler *BaseHandler)unifiedProcess() {
	log.WithField("url", handler.Request.URL.Path).
		WithField("time", time.Now().Format("2006-01-02 15:04:05.999")).
		Warn("404 this method is not implemented")
	handler.Response.Write(utils.Str2bytes("404 not found"))
}

// 用于处理http的get方法 args储存URL正则匹配之后的分组值 如果没有则为nil
func (handler *BaseHandler)Get() error {
	// todo::
	handler.unifiedProcess()
	return nil
}

func (handler *BaseHandler)Post() error {
	// todo::
	handler.unifiedProcess()
	//w.Write()
	return nil
}

func (handler *BaseHandler)Put() error {
	// todo::
	handler.unifiedProcess()
	return nil
}

func (handler *BaseHandler)Delete() error {
	// todo::
	handler.unifiedProcess()
	return nil
}

func (handler *BaseHandler)Patch() error {
	// todo::
	handler.unifiedProcess()
	return nil
}


func (handler *BaseHandler)Finish() error{
	handler.Response.Write(handler.buf.Bytes())
	return nil
}


// 一旦调用此方法 后边的写入内容不在写到response body中
func (handler *BaseHandler)OnFinish() {
	handler.isfinish = true
}

// 对请求体的参数进行解析 获取请求参数
func (handler *BaseHandler)GetArguments() map[string]interface{} {
	// go的标准包好像只对Content-Type为application/x-www-form-urlencoded的
	// 进行了处理 处理后 handler.Request.Body的内容就不存在了
	handler.Request.ParseForm()
	contentT := handler.Request.Header.Get("Content-Type")

	if contentT == "application/x-www-form-urlencoded" {
		if handler.Arguments == nil{
			handler.Arguments = make(map[string]interface{})
		}
		switch handler.Request.Method {
		case "POST":
			log.Debug("POST or PUT method")
			for k , v := range handler.Request.PostForm{
				handler.Arguments[k] = strings.Join(v, "")
			}
		case "PUT":
			for k , v := range handler.Request.PostForm{
				handler.Arguments[k] = strings.Join(v, "")
			}
		case "GET":
			for k , v := range handler.Request.Form{
				handler.Arguments[k] = strings.Join(v, "")
			}
		case "DELETE":
			for k , v := range handler.Request.Form{
				handler.Arguments[k] = strings.Join(v, "")
			}
		}
	}else if contentT == "application/json" {
		var buf bytes.Buffer
		buf.ReadFrom(handler.Request.Body)
		js, err := json.NewJson(buf.Bytes())
		if err != nil {
			log.Warn("paese form data fail", err)
		}else{
			handler.Arguments, _ = js.Map()
		}
	}
	return handler.Arguments
}

func (handler *BaseHandler)Write(content []byte) {
	if !handler.isfinish{
		handler.buf.Write(content)
	}
}
