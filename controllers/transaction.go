package controllers

import (
	"github.com/astaxie/beego"
	"hkzf/models"
)

type TransactionControoller struct {
	beego.Controller
}
func (this *TransactionControoller) SetValue()  {
	// key :订单编号:期数
	// value:From:To:金额:是否逾期:类型:备注

	orderId := this.GetString("orderId")
	issue := this.GetString("issue")

	if orderId == "" || issue == "" {
		handleResponse(this.Ctx.ResponseWriter, 400, "Request parameter orderId(issue) can't be empty")
		return
	}

	from := this.GetString("from")
	to := this.GetString("to")
	rent := this.GetString("rent")
	overdue := this.GetString("overdue")
	types := this.GetString("types")
	desc := this.GetString("desc")

	var(
		channelId=beego.AppConfig.String("channel_id_union")
		chainCodeId=beego.AppConfig.String("chaincode_id_transcation")
	)
	ccs, err := models.Initialize(channelId, chainCodeId, userId,"CORE_TRAN_CONFIG_FILE")
	if err !=nil{
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}
	defer ccs.Close()
	// chaincode查询需要的参数组装
	var args [][]byte
	args = append(args, []byte(orderId))
	args = append(args, []byte(issue))
	args = append(args, []byte(from))
	args = append(args, []byte(to))
	args = append(args, []byte(rent))
	args = append(args, []byte(overdue))
	args = append(args, []byte(types))
	args = append(args, []byte(desc))

	data, err := ccs.ChainCodeUpdate("set", args)

	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter, 200, data)

}


/*
	读数据的方法
 */
func (this *TransactionControoller) GetValue() {
	orderId := this.GetString("orderId")
	issue := this.GetString("issue")

	if orderId == "" || issue == "" {
		handleResponse(this.Ctx.ResponseWriter, 400, "Request parameter orderId(issue) can't be empty")
		return
	}

	// 依据key获取区块中记录的value信息,需要通过models中操作fabric-sdk-go,完成信息的读取
	var(
		channelId=beego.AppConfig.String("channel_id_union")
		chainCodeId=beego.AppConfig.String("chaincode_id_transcation")
	)
	ccs, err := models.Initialize(channelId, chainCodeId, userId,"CORE_TRAN_CONFIG_FILE")
	if err !=nil{
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}
	defer ccs.Close()
	// chaincode查询需要的参数组装
	var args [][]byte
	args = append(args, []byte(orderId))
	args = append(args, []byte(issue))

	response, err := ccs.ChainCodeQuery("get", args)

	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter, 200, response)
}
