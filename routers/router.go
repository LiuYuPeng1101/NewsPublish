package routers

import (
	"shanghaiyiqi/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//路由器过滤函数
	beego.InsertFilter("/article/*",beego.BeforeExec,filter)
    beego.Router("/", &controllers.MainController{},"get:ShowGet")
    //注册
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandlePost")
	//登录
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //登录列表页
    beego.Router("/article/ArticleList",&controllers.ArticleController{},"get:ShowArticleList")
    //添加文章
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
    //查看详情
    beego.Router("/article/ArticleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
    //文章编辑
    beego.Router("/article/ArticleUpdate",&controllers.ArticleController{},"get:ShowArticleUpdate;post:HandleArticleUpdate")
    //删除文章
    beego.Router("/article/DeleteArticle",&controllers.ArticleController{},"get:DeleteArticle")
    //添加分类
    beego.Router("/article/AddType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    //退出操作
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
	//删除类型
	beego.Router("/article/DeleteType",&controllers.ArticleController{},"get:DeleteType")
}

var filter = func(ctx *context.Context){
	userName:=ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302,"/login")
		return
	}
}
