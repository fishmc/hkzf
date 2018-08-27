package controllers

import (
	"github.com/astaxie/beego"
	"hkzf/models"
	"path"
	"os"
	"encoding/csv"
	"strconv"
	"io"
	"strings"
	"time"
	"fmt"
	"github.com/astaxie/beego/toolbox"
)

type CertificationController struct {
	beego.Controller
}
func (this *CertificationController)Check()  {

	//后端房屋认证
	houseId := this.GetString("houseId")
	id := this.GetString("id")

	if houseId == "" || id == ""{
		handleResponse(this.Ctx.ResponseWriter,400,"Reqest parameter houseId(or id) can't be empty")
		return
	}

	beego.Info(houseId + ":" + id)
	var (
		channelId = beego.AppConfig.String("channel_id_fgj")
		chainCodeId = beego.AppConfig.String("chaincode_id_house")
	)

	ccs, err := models.Initialize(channelId, chainCodeId, userId, "CORE_OFGJ_CONFIG_FILE")
	if err != nil{
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer  ccs.Close()

	var args [][]byte
	args = append(args,[]byte(houseId))
	args = append(args,[]byte(id))
	response, err := ccs.ChainCodeQuery("check", args)
	if err != nil{
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	handleResponse(this.Ctx.ResponseWriter,200,response)
}


func (this *CertificationController) RecordHouse() {
	var key = "house"
	// 用于上传房屋记录信息
	_, header, err := this.GetFile(key)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 400, err.Error())
		return
	}

	fileName := header.Filename
	beego.Info("文件名称:" + fileName)

	err = this.SaveToFile(key, path.Join("static/upload", fileName))
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}

	// 开启任务
	var myTask = "tk1"
	t := time.Now().Add(5 * time.Second)
	second := t.Second()
	minute := t.Minute()
	hour := t.Hour()
	spec := fmt.Sprintf("%d %d %d * * *", second, minute, hour)
	tk := toolbox.NewTask(myTask, spec, func() error {
		defer toolbox.StopTask()
		return myTask1(fileName)
	})

	toolbox.AddTask(myTask, tk)
	toolbox.StartTask()

	handleResponse(this.Ctx.ResponseWriter, 200, "ok")
}

func myTask1(fileName string) error {
	var (
		channelId   = beego.AppConfig.String("channel_id_fgj")
		chainCodeId = beego.AppConfig.String("chaincode_id_house")
	)

	ccs, err := models.Initialize(channelId, chainCodeId, userId, "CORE_OFGJ_CONFIG_FILE")
	if err != nil {
		beego.Error(err.Error())
		return err
	}
	defer ccs.Close()

	file, _ := os.Open(path.Join("static/upload", fileName))
	reader := csv.NewReader(file)

	var line = 0
	var lines []string

	for {
		line += 1
		linestr := strconv.Itoa(line)

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err !=nil{
			beego.Error(err.Error())
			lines=append(lines,linestr)
			continue
		}

		if len(record)!=3{
			beego.Error(err.Error())
			lines=append(lines,linestr)
			continue
		}

		var args [][]byte
		for _, str := range record {
			// 参数
			args = append(args, []byte(str))
		}
		_, err = ccs.ChainCodeUpdate("add", args)
		if err != nil {
			beego.Error(err.Error())
			lines=append(lines,linestr)
		}
	}

	if len(lines)>0{
		beego.Error("Error line:",strings.Join(lines,","))
	}else{
		beego.Info("write data success")
	}
	return nil
}
