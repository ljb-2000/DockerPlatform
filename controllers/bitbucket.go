package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

type BitbucketController struct {
	baseControllers
}

func (this *BitbucketController) Get() {
	username := this.GetSession("username").(string)

	bitbucketapiurl := beego.AppConfig.String("bitbucketapiurl")
	bitbucketuser := username
	bitbucketpass := GetPasswd(bitbucketuser)

	url := bitbucketapiurl + "/projects?limit=100"
	result, _ := RequestForAuth("GET", url, bitbucketuser, bitbucketpass, nil)

	var bitbucketlist interface{}
	json.Unmarshal(result, &bitbucketlist)

	this.Data["bitbucketlist"] = bitbucketlist
	this.Layout = "layout.html"
	this.TplName = "bitbucket.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func (this *BitbucketController) BitbucketRepos() {
	bitbucketapiurl := beego.AppConfig.String("bitbucketapiurl")
	bitbucketuser := this.GetSession("username").(string)
	bitbucketpass := GetPasswd(bitbucketuser)

	project := this.GetString("project")
	url := bitbucketapiurl + "/projects/" + project + "/repos"
	result, _ := RequestForAuth("GET", url, bitbucketuser, bitbucketpass, nil)

	var bitbucketrepos interface{}
	json.Unmarshal(result, &bitbucketrepos)

	this.Data["bitbucketrepos"] = bitbucketrepos
	this.Layout = "layout.html"
	this.TplName = "bitbucketrepos.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
