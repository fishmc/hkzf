package routers

import (
	"hkzf/controllers"
	"github.com/astaxie/beego"
)

func init() {

	//个人认证
    beego.Router("/auth", &controllers.AuthController{},"get:Check")
    beego.Router("/auth", &controllers.AuthController{},"post:RecordAuth")

    //房东认证
    beego.Router("/house",&controllers.CertificationController{},"get:Check")
    beego.Router("/house",&controllers.CertificationController{},"post:RecordHouse")

    //合同交易
    beego.Router("/contract",&controllers.ContractController{},"post:SetValue")
    beego.Router("/contract",&controllers.ContractController{},"get:GetValue")

}
