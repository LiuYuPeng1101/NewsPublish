package models

//model模式放的一般都是表的设计,增删改查全部都放到controllers下

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//创建结构体,其实就是创建表,在创建结构体的时候在orm方式中避免使用__下划线
//用户表
type User struct {
	Id int
	Name string
	Password string

	Article []*Article `orm:"reverse(many)"`
}
//文章表
type Article struct {
	Id int `orm:"pk;auto"`
	ArtiName string `orm:"size(20)"`
	Atime time.Time `orm:"auto_now"`
	Acount int `orm:"default(0);null"`
	Acontent string `orm:"size(500)"`
	Aimg string  `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk);on_delete(set_null);null"`
	User []*User `orm:"rel(m2m)"`
}
//类型表
type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`

	Article []*Article `orm:"reverse(many)"`
}

func init() {
//orm操作数据库
//1.连接数据库
err:=orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
if err !=nil{
	beego.Info("连接数据库失败")
}
//2.创建表
orm.RegisterModel(new(User),new(Article),new(ArticleType))
//3.生成表
//第一个参数表示给数据库起别名
//第二个参数表示是否强制更新
//第三个参数表示创建表的过程是否可见
orm.RunSyncdb("default",false,true)

}
