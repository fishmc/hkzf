package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)


/*
处理数据回复
*/

var userId = beego.AppConfig.String("user_id")

func handleResponse(response *context.Response,code int ,msg interface{})  {
	if code >= 400{
		beego.Error(msg)
	}else{
		beego.Info(msg)
	}
	response.WriteHeader(code)
	// 讲msg值进行一个类型的转换,如果转换成功ok(true)
	b , ok := msg.([]byte)
	if ok{
		response.Write(b)
	}else {
		s := msg.(string)
		response.Write([]byte(s))
	}
}
