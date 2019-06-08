package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["data"] = "china"
	c.TplName = "test.html"
}
func (c *MainController) ShowGet(){
	//获取orm对象
	//o:=orm.NewOrm()
	//插入数据
	//var user models.User
	//user.Name = "yupeng"
	//user.Password = "123456"
	//count,err:=o.Insert(&user)
	//if err !=nil{
	//	beego.Error("插入数据失败",err)
	//}
	//beego.Info("插入条数",count)
	//查询操作
	//var user models.User
	//user.Id = 1
	//err:=o.Read(&user,"Id")
	//if err !=nil{
	//	beego.Error("查询失败",err)
	//}
	//beego.Info("内容",user)
	//更新操作
	//var  user models.User
	//user.Id = 1
	//err:=o.Read(&user,"Id")
	//if err !=nil{
	//	beego.Error("查询失败",err)
	//}
	//user.Name = "liuyupeng"
	//count,err:=o.Update(&user)
	//if  err !=nil{
	//	beego.Error("更新失败",err)
	//}
	//beego.Info("更新条数",count)
	//删除操作
	//var user models.User
	//user.Id = 1
	//count,err:=o.Delete(&user)
	//if err !=nil{
	//	beego.Error("删除出错",err)
	//}
	//beego.Info("删除条数",count)

	//c.Data["data"] = "china"
	//c.TplName = "test.html"
}
