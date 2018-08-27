package controllers

import (
	"github.com/astaxie/beego"
	"hkzf/models"
	"path"
	"time"
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"os"
	"encoding/csv"
	"strconv"
	"io"
	"strings"
)

type AuthController struct {
	beego.Controller
}




func(this *AuthController)Check()  {
	// 发送数据:用户姓名和用户身份证号码
	name := this.GetString("name")
	id := this.GetString("id")
	if name == "" || id == ""{
		handleResponse(this.Ctx.ResponseWriter,400,"Request parameter name(or id) can't be empty")
		return
	}
	beego.Info(name + ":" + id)
	// 利用models当中的ChainCodeQuery查询当前用户的匹配结果和是否有不良个人记录的结果
	var(
		channelID = beego.AppConfig.String("channel_id_gaj")
		chainCodeID = beego.AppConfig.String("chaincode_id_auth")
		//userId      = beego.AppConfig.String("user_id")
	)
	css, err := models.Initialize(channelID, chainCodeID,userId,"CORE_OGAJ_CONFIG_FILE")

	if err != nil{
		handleResponse(this.Ctx.ResponseWriter ,500,err.Error())
		return
	}
	defer  css.Close()
	//args := [][]byte{[]byte(name), []byte(id)}
	args := [][]byte{[]byte(name),[]byte(id)}
	response, err := css.ChainCodeQuery("check", args)
	if err !=nil{
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	handleResponse(this.Ctx.ResponseWriter,200,response)
}











func (this *AuthController)RecordAuth()  {

	// 接收上传的文件
	// 保存文件:将文件保存在static/upload目录下面
	// 回复用户已经接收上传的文件
	// 开启任务:toolbox.NewTask()
	// 读文件
	// 将读取的每条记录信息都写到区块中

	// 定义读取文件的key
	var key string = "auth"

	beego.Info("receive file")
	// auth是一个key值,需要与前端约定
	_, header, err := this.GetFile(key)

	if err != nil {
		// 如果出错,没有读到文件 file和header都是nil
		handleResponse(this.Ctx.ResponseWriter, 400, err.Error())
		return
	}

	// 正常的操作情况下文件应该关闭,由于SaveToFile中含有了Close方法,下面的代码可以注释掉
	// defer file.Close()

	// 获取文件名
	fileName := header.Filename
	beego.Info("文件名称:", fileName)

	// 关于保存文件的路径:static前不能添加/
	err = this.SaveToFile(key, path.Join("static/upload", fileName))

	if err != nil {
		handleResponse(this.Ctx.ResponseWriter, 500, err.Error())
		return
	}

	// 开启任务
	// 指定任务的开启时间,保存完文件后的5秒钟执行写数据的任务
	// 任务是放到容器中进行管理,容器的开启和关闭

	var taskName string = "t1"

	// 当前时间+5秒  "* * * * * *"
	t := time.Now().Add(5 * time.Second)
	second := t.Second()
	minute := t.Minute()
	hour := t.Hour()
	spec := fmt.Sprintf("%d %d %d * * *", second, minute, hour)

	task := toolbox.NewTask(taskName, spec, func() error {
		beego.Info("start task")
		// 当我们任务执行完成后,停止
		defer toolbox.StopTask()
		return myTask(fileName)
	})
	// 注意:
	// task.Run() 立即执行

	// 将任务添加到容器
	toolbox.AddTask(taskName, task)
	// 开启任务执行
	toolbox.StartTask()

	handleResponse(this.Ctx.ResponseWriter, 200, "ok")

}







//耗时操作
// 耗时操作
func myTask(fileName string) error {
	// 读文件
	// 写数据到区块
	var (
		channelId   = beego.AppConfig.String("channel_id_gaj")
		chainCodeId = beego.AppConfig.String("chaincode_id_auth")
		//userId      = beego.AppConfig.String("user_id")
	)

	ccs, e := models.Initialize(channelId, chainCodeId, userId,"CORE_OGAJ_CONFIG_FILE")
	defer ccs.Close()

	if e != nil {
		beego.Error(e.Error())
		return e
	}

	file, _ := os.Open(path.Join("static/upload", fileName))
	reader := csv.NewReader(file)

	// 并没有中止,原因:可能某行数据有问题,当并不是所有的都有问题,只需要记录有问题的行信息就可以
	// 还有一种可能涉及异常,某行数据的写区块出错
	// 记录行数:可能会有多个行数据出问题,不会使用行的字符串拼接方式
	// 会定义一个字符串的数组,最终生成字符串时可以使用","进行分割行数

	var line = 0
	var lines []string

	for {
		line += 1
		linestr := strconv.Itoa(line)
		record, err := reader.Read()
		if err == io.EOF {
			// 文件结尾
			break
		}

		if err != nil {
			// 有异常需要处理
			lines = append(lines, linestr)
			continue
		}

		if len(record) != 3 {
			// 有异常需要处理
			lines = append(lines, linestr)
			continue
		}

		var args [][]byte

		for _, str := range record {
			//fmt.Print(str,"\t")
			args = append(args, []byte(str))
		}
		//fmt.Println()
		_, err = ccs.ChainCodeUpdate("add", args)

		if err != nil {
			// 有异常需要处理
			lines = append(lines, linestr)
		}
	}

	if len(lines) > 0 {
		beego.Error("Error lines:", strings.Join(lines, ","))
	}else{
		// 执行的一切顺利
		beego.Info("write data success")
	}

	return nil
}



