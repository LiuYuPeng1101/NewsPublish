package main

import (
	_ "shanghaiyiqi/routers"
	"github.com/astaxie/beego"
	_"shanghaiyiqi/models"
)

func main() {
	//连接视图函数和html
	beego.AddFuncMap("PrePage",ShowPrePage)
	beego.AddFuncMap("NextPage",ShowNextPage)
	beego.Run()
}

//视图函数：用于处理页面上一页和下一页的功能
func ShowPrePage(PageIndex int)int{
	if PageIndex == 1{
		return PageIndex
	}
	return PageIndex - 1
}
func ShowNextPage(PageIndex int,PageCount int)int{
	if PageIndex == PageCount{
		return PageCount
	}
	return PageIndex + 1
}
/*
视图函数:一般主要作为处理视图中简单的业务逻辑
1.创建视图函数
2.在视图中定义函数名
3.在beego.run()之前关联起来
 */


