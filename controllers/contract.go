package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"crypto/sha256"
	"io"
	"encoding/hex"
	"hkzf/models"
)

type ContractController struct {
	beego.Controller
}


func (this *ContractController) SetValue() {
	// 接收合同图片
	// 对图片进行SAH256处理--value
	// 约定:图片的名称为key

	file, header, err := this.GetFile("contract")
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 400, err.Error())
		return
	}

	// 获取订单编号
	fileName := header.Filename
	split := strings.Split(fileName, ".")
	key := split[0]

	beego.Info(key)

	// 获取图片的SHA256信息
	hash := sha256.New()
	// hash.Write([]byte)
	// 文件操作
	_, err = io.Copy(hash, file)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 400, err.Error())
		return
	}

	sum := hash.Sum(nil)
	value := hex.EncodeToString(sum)
	defer file.Close()

	// 写数据到区块
	// 写数据
	//var data []byte
	var(
		channelId=beego.AppConfig.String("channel_id_union")
		chainCodeId=beego.AppConfig.String("chaincode_id_contract")
	)
	ccs, err := models.Initialize(channelId, chainCodeId, userId,"CORE_CONT_CONFIG_FILE")
	if err !=nil{
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}
	defer ccs.Close()

	// 传递参数
	var args [][]byte
	args=append(args, []byte(key))
	args=append(args, []byte(value))

	beego.Info("key:"+key)
	beego.Info("value:"+value)

	response, err := ccs.ChainCodeUpdate("set", args)

	if err !=nil{
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}
	// 写数据成功的回复
	handleResponse(this.Ctx.ResponseWriter, 200, response)
}


/*
	读数据的方法
 */
func (this *ContractController) GetValue() {
	key := this.GetString("contractId")
	if key == "" {
		handleResponse(this.Ctx.ResponseWriter, 400, "Request parameter contractId can't be empty")
		return
	}
	beego.Info("key:", key)

	// 依据key获取区块中记录的value信息,需要通过models中操作fabric-sdk-go,完成信息的读取
	var(
		channelId=beego.AppConfig.String("channel_id_union")
		chainCodeId=beego.AppConfig.String("chaincode_id_contract")
	)
	ccs, err := models.Initialize(channelId, chainCodeId, userId,"CORE_CONT_CONFIG_FILE")
	if err !=nil{
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}
	defer ccs.Close()
	// chaincode查询需要的参数组装
	var args [][]byte
	args = append(args, []byte(key))

	response, err := ccs.ChainCodeQuery("get", args)

	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter, 200, response)
}