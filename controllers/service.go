package controllers

import (
	"encoding/json"
	// "github.com/astaxie/beego"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type ServiceController struct {
	baseControllers
}

func (this *ServiceController) Get() {
	username := this.GetSession("username").(string)

	namespaces := ""
	if CheckAdmin(username) {
		namespaces = this.GetString("namespaces")
	} else {
		namespaces = username
	}

	url := "http://10.10.7.175:8081/services/list"

	//json 序列化
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	//post数据
	result, _ := Request("POST", url, jsonStr)

	var servicelist []v1.Service
	json.Unmarshal(result, &servicelist)

	this.Data["servicelist"] = servicelist
	this.Layout = "layout.html"
	this.TplName = "service.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
