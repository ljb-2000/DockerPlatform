package controllers

import (
	"encoding/json"
	// "github.com/astaxie/beego"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type IngressController struct {
	baseControllers
}

func (this *IngressController) Get() {
	username := this.GetSession("username").(string)

	namespaces := ""
	if CheckAdmin(username) {
		namespaces = this.GetString("namespaces")
	} else {
		namespaces = username
	}

	url := "http://10.10.7.175:8081/ingress/list"

	//json 序列化
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	//post数据
	result, _ := Request("POST", url, jsonStr)

	var ingresslist []v1beta1.Ingress
	json.Unmarshal([]byte(result), &ingresslist)

	this.Data["ingresslist"] = ingresslist

	this.Layout = "layout.html"
	this.TplName = "ingress.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
