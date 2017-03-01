/*
* 时间: 2017年03月01日13:12:19
* 说明: 一个简单的web框架 仿照tornado的框架模型
* 作者: wupeaking
*/

package web

import (
	"net/http"
	"github.com/wupeaking/typhon/app/utils"
	log "github.com/Sirupsen/logrus"
	"time"
)

// 一个基本的handler应该至少实现http的PUT DELET POST GET PATCH五种方法
type BaseHandler struct {

}

func (handler *BaseHandler)Get(w http.ResponseWriter, r *http.Request) error {
	// todo::
	log.WithField("url", r.URL.Path).
		WithField("time", time.Now().Format("2006-01-02 15:04:05.999")).
		Warn("404 this method is not implemented")
	w.Write(utils.Str2bytes("404 not found"))
	return nil
}

func (handler *BaseHandler)Post(w http.ResponseWriter, r *http.Request) error {
	// todo::
	//w.Write()
	return nil
}

func (handler *BaseHandler)Put(w http.ResponseWriter, r *http.Request) error {
	// todo::
	//w.Write()
	return nil
}

func (handler *BaseHandler)Delete(w http.ResponseWriter, r *http.Request) error {
	// todo::
	//w.Write()
	return nil
}

func (handler *BaseHandler)Patch(w http.ResponseWriter, r *http.Request) error {
	// todo::
	//w.Write()
	return nil
}

