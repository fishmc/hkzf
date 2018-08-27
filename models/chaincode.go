package models

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/astaxie/beego"
)

// 个人认证
// 发送姓名和身份证信息，回复：是否通过认证，是否有不良个人记录


type ChaincodeSpec struct {
	client apitxn.ChannelClient
	chainCodeID string
}
func Initialize(channelID ,chainCodeID,userId,conf string) (*ChaincodeSpec,error) {

	config := beego.AppConfig.String("conf")
	sdk, err := getSDK(config)
	if err != nil{
		return nil,err
	}
	client, err := sdk.NewChannelClient(channelID, userId)
	if err != nil{
		return nil ,err
	}
	return  &ChaincodeSpec{client,channelID},nil
}
func (this *ChaincodeSpec)ChainCodeUpdate(function string,args [][]byte) (response []byte,err error) {
	request := apitxn.ExecuteTxRequest{ChaincodeID:this.chainCodeID,Fcn:function,Args:args}
	id, err := this.client.ExecuteTx(request)
	return []byte(id.ID),nil
}
func (this *ChaincodeSpec)ChainCodeQuery(function string,args [][]byte) (response []byte,err error) {
	request := apitxn.QueryRequest{this.chainCodeID,function,args}
	return this.client.Query(request)
}

func (this *ChaincodeSpec)Close()  {
	this.client.Close()
}

func getSDK(config string) (*fabapi.FabricSDK,error) {
	options := fabapi.Options{ConfigFile:config}
	sdk, err := fabapi.NewSDK(options)
	if err != nil{
		beego.Error(err.Error())
	}
	return sdk,err
}

