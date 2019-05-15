package models
/*

import (
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//连接数据库
	conn,err:=sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err !=nil{
		beego.Error("数据库连接错误",err)
		return
	}
	defer conn.Close()

	//创建表
	//res,err:=conn.Exec("create table user(name varchar(40),password varchar(40));")
	//if err !=nil{
	//	beego.Error("创建表错误",err)
	//	beego.Info("创建表信息",res)
	//	return
	//}
	//增删该查--增
	//res,err:=conn.Exec("insert into user(name,password) values (?,?)","yupeng","123456")
	//if err !=nil{
	//	beego.Error("增加数据错误",err)
	//	beego.Info("增加数据信息",res)
	//	return
	//}
	//查询
	res,err:=conn.Query("select name from user")
	var name string
	for res.Next(){
		res.Scan(&name)
		beego.Info(res)
	}
}
*/
