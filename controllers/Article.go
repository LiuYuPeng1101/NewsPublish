package controllers

import (
	"bytes"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"path"
	"shanghaiyiqi/models"
	"time"
)

type ArticleController struct {
	beego.Controller
}

func (this*ArticleController)ShowArticleList(){
	userName:=this.GetSession("userName")
	if userName ==nil{
		this.Redirect("/login",302)
		return
	}
	//从数据库获取数据
	//这里会运用关于orm的高级查询
	o:=orm.NewOrm()
	qs:=o.QueryTable("article")
	var Articles []models.Article
	//_,err:=qs.All(&Articles)
	//if err !=nil{
	//	beego.Info("查询数据错误")
	//}
	//获取总记录数
	TypeName:=this.GetString("select")
	beego.Info("TypeName",TypeName)
	var Count int64
	//获取总页数
	PageSize  := 2
	//从index.html获取页码
	PageIndex,err:=this.GetInt("PageIndex")
	if err !=nil{
		PageIndex = 1
	}
	//确定数据的起始位置
	start:=(PageIndex - 1) * PageSize
	if TypeName == ""{
		Count,_=qs.Count()
		beego.Info("Count",Count)
	}else{
		Count,_=qs.Limit(PageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",TypeName).Count()
	}
	PageCount := math.Ceil(float64(Count) / float64(PageSize))
	//获取数据库文章类型列表
	var ArticleTypes []models.ArticleType
	//需要把文章类型的数据存储到redis中,redis通常存储的就是经常访问但又不改动的数据，在这个页面中，文章类型就属于这个数据
	//序列化和反序列化
	conn,err:=redis.Dial("tcp","172.16.107.175:6379")
	if err !=nil{
		beego.Info("redis数据库连接失败")
	}
	defer conn.Close()
	res,err:=conn.Do("get","ArticleTypes")
	if err !=nil{
		o.QueryTable("ArticleType").All(&ArticleTypes)
		var buffer bytes.Buffer
		enc:=gob.NewEncoder(&buffer) //编码器
		enc.Encode(ArticleTypes) //编码
		conn.Do("set","ArticleTypes",buffer.Bytes())
	}else{
		data,_:=redis.Bytes(res,err)
		dec:=gob.NewDecoder(bytes.NewReader(data))
		dec.Decode(&ArticleTypes)
	}
	//获取前端文章类型
	if TypeName == ""{
		qs.Limit(PageSize,start).RelatedSel("ArticleType").All(&Articles)
	}else{
		qs.Limit(PageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",TypeName).All(&Articles)
	}
	//传递数据
	this.Data["TypeName"] = TypeName
	this.Data["PageIndex"] = PageIndex
	this.Data["Count"] = Count
	this.Data["PageCount"] = int(PageCount)
	this.Data["Articles"] = Articles
	this.Data["ArticleTypes"] = ArticleTypes
	this.Data["userName"] = userName

	this.Layout="layout.html"
	this.TplName = "index.html"
}

func (this*ArticleController)ShowAddArticle(){
	//获取数据
	o:=orm.NewOrm()
	var ArticleType []models.ArticleType
	o.QueryTable("ArticleType").All(&ArticleType)
	this.Data["ArticleType"] = ArticleType
	userName:=this.GetSession("userName")
	this.Data["userName"] = userName
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

func (this*ArticleController)ShowArticleDetail(){
	//获取数据：从html中拿到文章id
	ContentId,err:=this.GetInt("ArticleDetail")
	//校验数据
	if err !=nil{
		beego.Info("获取id失败")
	}
	o:=orm.NewOrm()
	var Article models.Article
	Article.Id = ContentId
	//o.Read(&Article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",ContentId).One(&Article)
	//修改阅读次数,每查看一次详情，增加一次
	Article.Acount +=1
	o.Update(&Article)
	//多对多插入
	m2m:=o.QueryM2M(&Article,"User")
	userName:=this.GetSession("userName")
	if userName == ""{
		this.Redirect("/login",302)
		return
	}
	var User models.User
	User.Name = userName.(string)
	o.Read(&User,"Name")
	m2m.Add(User)
	//查询
	var users []models.User
	o.QueryTable("User").Filter("Article__Article__Id",ContentId).Distinct().All(&users)

	//返回数据
	this.Data["Article"] = Article
	this.Data["users"] = users
	userlayout:=this.GetSession("userName")
	this.Data["userName"] = userlayout
	this.Layout="layout.html"
	this.TplName="content.html"
}

func (this*ArticleController)ShowArticleUpdate(){
	//获取数据
	ContentId,err:=this.GetInt("ArticleUpdate")
	//校验数据
	if err !=nil{
		beego.Info("文章id不正确")
	}
	//操作数据
	o:=orm.NewOrm()
	var Article models.Article
	Article.Id = ContentId
	o.Read(&Article)
	//返回数据
	this.Data["Article"] = Article
	userName:=this.GetSession("userName")
	this.Data["userName"] = userName
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

func (this*ArticleController)HandleAddArticle(){
	//获取数据
	ArticleName:=this.GetString("articleName")
	Content:=this.GetString("content")
	//校验数据
	if ArticleName == "" || Content == ""{
		this.Data["errmsg"] = "添加数据不完整"
		this.TplName = "add.html"
		return
	}
	//图片上传功能实现
	//图片上传我们需要考虑三点问题,图片大小,图片格式,避免重命名
	file,head,err:=this.GetFile("uploadname")
	defer file.Close()
	if err !=nil{
		this.Data["errmsg"] = "上传图片错误,请重新上传"
		this.TplName = "add.html"
		return
	}
	//图片大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "上传图片太大,请重新上传"
		this.TplName = "add.html"
		return
	}
	//图片格式
	fileExt:=path.Ext(head.Filename)
	if fileExt !=".png" && fileExt !=".jpg"{
		this.Data["errmsg"] = "上传图片格式不正确,请重新上传"
		this.TplName = "add.html"
		return
	}
	//避免重命名
	fileName:=time.Now().Format("2006-01-02-15-04-05")
	//保存文件
	this.SaveToFile("uploadname","./static/img"+fileName+fileExt)
	//操作数据
	o:=orm.NewOrm()
	var Article models.Article
	Article.Acontent = Content
	Article.ArtiName = ArticleName
	Article.Aimg = "/static/img"+fileName+fileExt
	//获取文章类型
	Type:=this.GetString("select")
	var ArticleType models.ArticleType
	ArticleType.TypeName = Type
	o.Read(&ArticleType,"TypeName")
	Article.ArticleType = &ArticleType
	o.Insert(&Article)
	//返回结果
	this.Redirect("/article/ArticleList",302)
}

//封装函数--函数功能：得到从前端传递的图片
func UploadFile(this*beego.Controller,FilePath string)string{
	//图片上传功能实现
	//图片上传我们需要考虑三点问题,图片大小,图片格式,避免重命名
	file,head,err:=this.GetFile(FilePath)
	if err !=nil{
		this.Data["errmsg"] = "上传图片错误,请重新上传"
		this.TplName = "add.html"
		return ""
	}
	if head.Filename == ""{
		return "NoImg"
	}
	defer file.Close()
	//图片大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "上传图片太大,请重新上传"
		this.TplName = "add.html"
		return ""
	}
	//图片格式
	fileExt:=path.Ext(head.Filename)
	if fileExt !=".png" && fileExt !=".jpg"{
		this.Data["errmsg"] = "上传图片格式不正确,请重新上传"
		this.TplName = "add.html"
		return ""
	}
	//避免重命名
	fileName:=time.Now().Format("2006-01-02-15-04-05")
	//保存文件
	this.SaveToFile(FilePath,"./static/img"+fileName+fileExt)
	return  "/static/img"+fileName+fileExt
}

func (this*ArticleController)HandleArticleUpdate(){
	//获取数据
	Id,err:=this.GetInt("ArticleUpdate")
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	FilePath:=UploadFile(&this.Controller,"uploadname")
	//校验数据
	if err !=nil || articleName == "" || content == ""{
		beego.Info("上传数据错误")
	}
	//操作数据
	o:=orm.NewOrm()
	var Article models.Article
	Article.Id = Id
	o.Read(&Article)

	Article.ArtiName = articleName
	Article.Acontent = content
	if FilePath != "NoImg"{
		Article.Aimg = FilePath
	}
	o.Update(&Article)
	//返回视图
	this.Redirect("/ArticleList",302)
}

func (this*ArticleController)DeleteArticle(){
	//获取数据
	Id,err:=this.GetInt("ArticleId")
	//校验数据
	if err !=nil{
		beego.Info("传递数据出错")
	}
	//操作数据
	o:=orm.NewOrm()
	var Article models.Article
	Article.Id = Id
	o.Delete(&Article)
	//返回数据
	this.Redirect("/ArticleList",302)




}

func (this*ArticleController)ShowAddType(){
	o:=orm.NewOrm()
	var ArticleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&ArticleTypes)
	//传递数据
	this.Data["ArticleTypes"] = ArticleTypes
	userName:=this.GetSession("userName")
	this.Data["userName"] = userName
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

func (this*ArticleController)HandleAddType(){
	//获取数据
	TypeName:=this.GetString("typeName")
	//校验数据
	if TypeName == ""{
		beego.Info("传递数据失败")
	}
	//操作数据
	o:=orm.NewOrm()
	var ArticleType models.ArticleType
	ArticleType.TypeName = TypeName
	o.Insert(&ArticleType)
	//返回视图
	this.Redirect("/article/AddType",302)
}

func (this*ArticleController)DeleteType(){
	Id,err:=this.GetInt("Id")
	if err !=nil{
		beego.Info("传递数据失败")
		return
	}
	o:=orm.NewOrm()
	var ArticleType models.ArticleType
	ArticleType.Id = Id
	o.Delete(&ArticleType)
	this.Redirect("/article/AddType",302)


}

