package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"shanghaiyiqi/models"
)

type UserController struct {
	beego.Controller
}

//显示注册页面
func (this *UserController)ShowRegister(){
	//指定视图
	this.TplName = "register.html"
}
//显示登录页面
func (this *UserController)ShowLogin(){
	userName:=this.Ctx.GetCookie("userName")
	if userName == ""{
		this.Data["userName"] = ""
		this.Data["check"] = ""
	}else{
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}

	//指定视图
	this.TplName = "login.html"
}
//处理注册数据
func (this *UserController)HandlePost(){
	//获取数据
	UserName:=this.GetString("userName")
	PassWord:=this.GetString("password")
	beego.Info(UserName,PassWord)
	//校验数据
	if UserName == "" || PassWord == ""{
		this.Data["errmsg"] = "注册数据不完整,请重新注册"
		this.TplName = "register.html"
		return
	}
	//操作数据
	o:=orm.NewOrm()
	var user models.User
	user.Name = UserName
	user.Password = PassWord
	o.Insert(&user)
	//返回页面
	//this.Ctx.WriteString("注册成功")
	this.Redirect("/login",302)
}
//处理登录数据
func (this*UserController)HandleLogin(){
	//获取数据
	UserName:=this.GetString("userName")
	PassWord:=this.GetString("password")
	beego.Info(UserName,PassWord)
	//校验数据
	if UserName == "" || PassWord == ""{
		this.Data["errmsg"] = "登录数据不完整,请重新登录"
		this.TplName = "login.html"
		return
	}
	//操作数据
	o:=orm.NewOrm()
	var user models.User
	user.Name = UserName
	err:=o.Read(&user,"name")
	if err !=nil{
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}
	if user.Password != PassWord{
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}
	//返回结果
	//this.Ctx.WriteString("登录成功")
	data:=this.GetString("remember")
	if data == "on"{
		this.Ctx.SetCookie("userName",UserName,100)
	}else{
		this.Ctx.SetCookie("userName",UserName,-1)
	}
	this.SetSession("userName",UserName)
	this.Redirect("/article/ArticleList",302)
}
//退出
func (this*UserController)Logout(){
	this.DelSession("userName")
	this.Redirect("/login",302)
}

