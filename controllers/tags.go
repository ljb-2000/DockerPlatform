package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

type TagsController struct {
	baseControllers
}

func (this *TagsController) Get() {
	username := this.GetSession("username").(string)
	repos := this.GetString("repos", "")
	projectid := this.GetString("projectid", "")
	harborhost := beego.AppConfig.String("harborhost")
	url := harborhost + "/api/repositories/" + repos + "/tags"

	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")
	result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)

	var f []interface{}
	json.Unmarshal(result, &f)

	tags := make(map[string]interface{})
	for _, v := range f {
		tags[v.(string)] = GetManifest(repos, v.(string))
	}

	this.Data["harbor"] = "harbor.gqichina.com"
	this.Data["projectid"] = projectid
	this.Data["repos"] = repos
	this.Data["tags"] = tags
	this.Data["namespaces"] = username
	this.Layout = "layout.html"
	this.TplName = "tags.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func GetManifest(repos, tags string) interface{} {
	harborhost := beego.AppConfig.String("harborhost")
	url := harborhost + "/api/repositories/" + repos + "/tags/" + tags + "/manifest"
	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")
	result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)

	var f interface{}
	json.Unmarshal(result, &f)

	m := f.(map[string]interface{})

	var m1 = []byte(m["config"].(string))
	var f1 interface{}
	json.Unmarshal(m1, &f1)
	return f1
}
